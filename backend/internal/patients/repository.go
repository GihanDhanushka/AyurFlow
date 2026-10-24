package patients

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func readPatient(ctx context.Context, tx pgx.Tx, id string) (Patient, error) {
	var patient Patient
	var phone pgtype.Text
	err := tx.QueryRow(ctx, `SELECT id::text, display_name, phone, record_version, created_at, updated_at
		FROM patients WHERE id = $1`, id).Scan(
		&patient.ID, &patient.DisplayName, &phone, &patient.Version,
		&patient.CreatedAt, &patient.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return Patient{}, ErrNotFound
	}
	if err != nil {
		return Patient{}, err
	}
	if phone.Valid {
		patient.Phone = &phone.String
	}
	var identity IdentityInput
	err = tx.QueryRow(ctx, `SELECT identity_type, original_value FROM patient_identities
		WHERE patient_id = $1`, id).Scan(&identity.Type, &identity.Value)
	if err == nil {
		patient.Identity = &identity
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return Patient{}, err
	}
	return patient, nil
}

func identityError(err error) error {
	var dbErr *pgconn.PgError
	if errors.As(err, &dbErr) && dbErr.Code == "23505" &&
		dbErr.ConstraintName == "patient_identities_type_canonical_unique" {
		return ErrDuplicateIdentity
	}
	return err
}
