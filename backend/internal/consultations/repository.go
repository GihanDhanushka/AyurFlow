package consultations

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func readDraft(ctx context.Context, tx pgx.Tx, patientID, id string) (Consultation, error) {
	var draft Consultation
	err := tx.QueryRow(ctx, `SELECT id::text, patient_id::text, status, draft_text, record_version,
		created_at, created_by::text, updated_at, updated_by::text
		FROM consultations WHERE id=$1 AND patient_id=$2`, id, patientID).Scan(
		&draft.ID, &draft.PatientID, &draft.Status, &draft.DraftText, &draft.Version,
		&draft.CreatedAt, &draft.CreatedBy, &draft.UpdatedAt, &draft.UpdatedBy,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Consultation{}, ErrNotFound
	}
	return draft, err
}
