//go:build integration

package consultations

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"ayurflow/backend/internal/testutil"

	"github.com/jackc/pgx/v5"
)

type assignedAccess struct {
	actor   string
	patient string
}

func (a *assignedAccess) require(actor, patient string) error {
	if a.actor != actor || a.patient != patient {
		return ErrForbidden
	}
	return nil
}

func (a *assignedAccess) CanCreate(_ context.Context, _ pgx.Tx, actor, patient string) error {
	return a.require(actor, patient)
}
func (a *assignedAccess) CanRead(_ context.Context, _ pgx.Tx, actor, patient, _ string) error {
	return a.require(actor, patient)
}
func (a *assignedAccess) CanEdit(_ context.Context, _ pgx.Tx, actor, patient, _ string) error {
	return a.require(actor, patient)
}

func TestConsultationDraftServiceIntegration(t *testing.T) {
	pool := testutil.Postgres(t)
	ctx, cancel := context.WithTimeout(t.Context(), time.Minute)
	defer cancel()
	var actor, other, patientID, otherPatientID string
	for i, target := range []*string{&actor, &other} {
		if err := pool.QueryRow(ctx, `INSERT INTO staff_accounts(login_identifier,password_hash)
			VALUES ($1, '$argon2id$synthetic-test-only') RETURNING id::text`,
			fmt.Sprintf("consultation-%d-%d", time.Now().UnixNano(), i)).Scan(target); err != nil {
			t.Fatal(err)
		}
	}
	for _, target := range []*string{&patientID, &otherPatientID} {
		if err := pool.QueryRow(ctx, `INSERT INTO patients(display_name,created_by,updated_by)
			VALUES ('Synthetic patient',$1,$1) RETURNING id::text`, actor).Scan(target); err != nil {
			t.Fatal(err)
		}
	}
	service := NewService(pool, &assignedAccess{actor: actor, patient: patientID})
	if _, err := service.Create(ctx, other, patientID, "Unauthorized note"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("other actor created consultation: %v", err)
	}
	if _, err := service.Create(ctx, actor, otherPatientID, "Wrong patient"); !errors.Is(err, ErrForbidden) {
		t.Fatalf("other patient consultation created: %v", err)
	}
	created, err := service.Create(ctx, actor, patientID, "සිංහල synthetic draft")
	if err != nil || created.Status != "draft" || created.Version != 1 || created.PatientID != patientID {
		t.Fatalf("draft without plan: %+v, %v", created, err)
	}
	if got, err := service.Get(ctx, actor, patientID, created.ID); err != nil || got.DraftText != created.DraftText {
		t.Fatalf("authorized draft read: %+v, %v", got, err)
	}
	if _, err := service.Get(ctx, other, patientID, created.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("other actor read clinical text: %v", err)
	}
	if _, err := service.Get(ctx, actor, patientID, otherPatientID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("unknown consultation must be hidden: %v", err)
	}

	// Concurrent edits with the same version serialize on the row; one is stale.
	var wg sync.WaitGroup
	results := make(chan error, 2)
	start := make(chan struct{})
	for _, note := range []string{"First edit", "Second edit"} {
		wg.Add(1)
		go func(note string) {
			defer wg.Done()
			<-start
			_, err := service.Update(ctx, actor, patientID, created.ID, created.Version, note)
			results <- err
		}(note)
	}
	close(start)
	wg.Wait()
	close(results)
	successes, stale := 0, 0
	for err := range results {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrStaleVersion):
			stale++
		default:
			t.Fatalf("edit race error: %v", err)
		}
	}
	if successes != 1 || stale != 1 {
		t.Fatalf("edit race: success=%d stale=%d", successes, stale)
	}
	current, err := service.Get(ctx, actor, patientID, created.ID)
	if err != nil || current.Version != 2 {
		t.Fatalf("committed draft after race: %+v, %v", current, err)
	}
	var auditCount, leakedBodies int
	if err := pool.QueryRow(ctx, `SELECT count(*), count(*) FILTER (WHERE reason IS NOT NULL)
		FROM audit_events WHERE object_type='consultation' AND object_id=$1`, created.ID).Scan(&auditCount, &leakedBodies); err != nil || auditCount != 2 || leakedBodies != 0 {
		t.Fatalf("audit metadata mismatch: count=%d reason-filled=%d err=%v", auditCount, leakedBodies, err)
	}
	// A failed audit insert must leave the committed draft and version intact.
	if _, err := pool.Exec(ctx, `CREATE FUNCTION reject_consultation_audit() RETURNS trigger LANGUAGE plpgsql AS $$
		BEGIN IF NEW.action='consultation.updated' THEN RAISE EXCEPTION 'synthetic audit failure'; END IF;
		RETURN NEW; END $$`); err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `CREATE TRIGGER reject_consultation_audit BEFORE INSERT ON audit_events
		FOR EACH ROW EXECUTE FUNCTION reject_consultation_audit()`); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Update(ctx, actor, patientID, created.ID, current.Version, "Must roll back"); err == nil {
		t.Fatal("audit failure did not abort draft edit")
	}
	after, err := service.Get(ctx, actor, patientID, created.ID)
	if err != nil || after.Version != current.Version || after.DraftText != current.DraftText {
		t.Fatalf("failed audit changed draft: %+v, %v", after, err)
	}
}
