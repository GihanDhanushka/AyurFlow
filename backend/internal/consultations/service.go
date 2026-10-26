package consultations

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"ayurflow/backend/internal/audit"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrForbidden      = errors.New("consultations: access denied")
	ErrNotFound       = errors.New("consultations: not found")
	ErrInvalidDraft   = errors.New("consultations: invalid draft")
	ErrStaleVersion   = errors.New("consultations: stale version")
	ErrNotDraft       = errors.New("consultations: no longer a draft")
)

// AccessPolicy must check an active clinician's patient assignment in the
// caller's transaction. No policy, including a nil policy, grants access.
// The assignment source is pending DEC02; protected routes stay unregistered.
type AccessPolicy interface {
	CanCreate(context.Context, pgx.Tx, string, string) error
	CanRead(context.Context, pgx.Tx, string, string, string) error
	CanEdit(context.Context, pgx.Tx, string, string, string) error
}

type Consultation struct {
	ID        string    `json:"id"`
	PatientID string    `json:"patientId"`
	Status    string    `json:"status"`
	DraftText string    `json:"draftText"`
	Version   int64     `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	CreatedBy string    `json:"createdBy"`
	UpdatedAt time.Time `json:"updatedAt"`
	UpdatedBy string    `json:"updatedBy"`
}

type Service struct {
	pool   *pgxpool.Pool
	access AccessPolicy
}

func NewService(pool *pgxpool.Pool, access AccessPolicy) *Service {
	return &Service{pool: pool, access: access}
}

func (s *Service) Create(ctx context.Context, actor, patientID, draftText string) (Consultation, error) {
	if !validUUID(actor) {
		return Consultation{}, ErrForbidden
	}
	if !validUUID(patientID) {
		return Consultation{}, ErrNotFound
	}
	if !validDraft(draftText) {
		return Consultation{}, ErrInvalidDraft
	}
	if s == nil || s.pool == nil || s.access == nil {
		return Consultation{}, ErrForbidden
	}
	var created Consultation
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := s.access.CanCreate(ctx, tx, actor, patientID); err != nil {
			return err
		}
		var exists bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM patients WHERE id=$1)`, patientID).Scan(&exists); err != nil {
			return err
		}
		if !exists {
			return ErrNotFound
		}
		var id string
		if err := tx.QueryRow(ctx, `INSERT INTO consultations(patient_id, draft_text, created_by, updated_by)
			VALUES ($1, $2, $3, $3) RETURNING id::text`, patientID, draftText, actor).Scan(&id); err != nil {
			return err
		}
		if err := audit.Write(ctx, tx, audit.Event{
			ActorID: actor, Action: "consultation.created",
			Object: audit.Object{Type: "consultation", ID: id},
		}); err != nil {
			return err
		}
		var err error
		created, err = readDraft(ctx, tx, patientID, id)
		return err
	})
	return created, err
}

func (s *Service) Get(ctx context.Context, actor, patientID, id string) (Consultation, error) {
	if !validUUID(actor) {
		return Consultation{}, ErrForbidden
	}
	if !validUUID(patientID) || !validUUID(id) {
		return Consultation{}, ErrNotFound
	}
	if s == nil || s.pool == nil || s.access == nil {
		return Consultation{}, ErrForbidden
	}
	var found Consultation
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := s.access.CanRead(ctx, tx, actor, patientID, id); err != nil {
			return err
		}
		var err error
		found, err = readDraft(ctx, tx, patientID, id)
		return err
	})
	return found, err
}

func (s *Service) Update(ctx context.Context, actor, patientID, id string, version int64, draftText string) (Consultation, error) {
	if !validUUID(actor) {
		return Consultation{}, ErrForbidden
	}
	if !validUUID(patientID) || !validUUID(id) {
		return Consultation{}, ErrNotFound
	}
	if version <= 0 || !validDraft(draftText) {
		return Consultation{}, ErrInvalidDraft
	}
	if s == nil || s.pool == nil || s.access == nil {
		return Consultation{}, ErrForbidden
	}
	var updated Consultation
	err := pgx.BeginFunc(ctx, s.pool, func(tx pgx.Tx) error {
		if err := s.access.CanEdit(ctx, tx, actor, patientID, id); err != nil {
			return err
		}
		var currentVersion int64
		var status string
		err := tx.QueryRow(ctx, `SELECT record_version, status FROM consultations
			WHERE id=$1 AND patient_id=$2 FOR UPDATE`, id, patientID).Scan(&currentVersion, &status)
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return err
		}
		if status != "draft" {
			return ErrNotDraft
		}
		if currentVersion != version {
			return ErrStaleVersion
		}
		if _, err := tx.Exec(ctx, `UPDATE consultations SET draft_text=$3,
			record_version=record_version+1, updated_at=clock_timestamp(), updated_by=$4
			WHERE id=$1 AND patient_id=$2`, id, patientID, draftText, actor); err != nil {
			return err
		}
		if err := audit.Write(ctx, tx, audit.Event{
			ActorID: actor, Action: "consultation.updated",
			Object: audit.Object{Type: "consultation", ID: id},
			ChangedFields: []string{"draft_text"},
		}); err != nil {
			return err
		}
		updated, err = readDraft(ctx, tx, patientID, id)
		return err
	})
	return updated, err
}

func validDraft(value string) bool {
	return utf8.ValidString(value) && len(value) <= 65536 && strings.TrimSpace(value) != ""
}

func validUUID(value string) bool {
	var id pgtype.UUID
	return id.Scan(value) == nil && id.Valid && id.Bytes != [16]byte{}
}
