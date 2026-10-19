// seedstaff creates an explicitly requested synthetic staff account for local development.
// It is not an account-management endpoint or a production bootstrap mechanism.
package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"ayurflow/backend/internal/auth"
	"ayurflow/backend/internal/config"
	"ayurflow/backend/internal/platform/database"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var loginPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{2,63}$`)

var seedRoles = map[string]bool{
	"administrator": true, "receptionist": true, "clinician": true,
	"therapist": true, "dispensary": true, "billing": true, "report_reader": true,
}

type input struct {
	login    string
	password string
	roles    []string
}

func main() {
	if err := run(); err != nil {
		// Do not print environment variables, database errors or account details.
		slog.Error("Staff seed failed; check development environment, input and migrations")
		os.Exit(1)
	}
	slog.Info("Synthetic staff seed complete")
}

func run() error {
	values, err := readInput(os.Getenv)
	if err != nil {
		return err
	}
	// Keep the secret out of child processes started after input validation.
	_ = os.Unsetenv("SEED_STAFF_PASSWORD")
	hash, err := auth.HashPassword(values.password)
	if err != nil {
		return err
	}
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	pool, err := database.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	return seed(ctx, pool, values.login, hash, values.roles)
}

func readInput(getenv func(string) string) (input, error) {
	if getenv("APP_ENV") != "development" {
		return input{}, errors.New("staff seed requires explicit development environment")
	}
	login := strings.ToLower(strings.TrimSpace(getenv("SEED_STAFF_LOGIN")))
	password := getenv("SEED_STAFF_PASSWORD")
	if !loginPattern.MatchString(login) || !utf8.ValidString(password) || utf8.RuneCountInString(password) < 12 || len(password) > 1024 {
		return input{}, errors.New("invalid staff seed input")
	}
	roleInput := getenv("SEED_STAFF_ROLES")
	if roleInput == "" {
		return input{}, errors.New("staff seed roles are required")
	}
	roles := strings.Split(roleInput, ",")
	seen := make(map[string]bool, len(roles))
	for i, role := range roles {
		role = strings.TrimSpace(role)
		if !seedRoles[role] || seen[role] {
			return input{}, errors.New("invalid staff seed role")
		}
		seen[role] = true
		roles[i] = role
	}
	return input{login: login, password: password, roles: roles}, nil
}

func seed(ctx context.Context, pool *pgxpool.Pool, login, passwordHash string, roles []string) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var staffID string
	err = tx.QueryRow(ctx, `INSERT INTO staff_accounts(login_identifier, password_hash)
		VALUES ($1, $2) ON CONFLICT (login_identifier) DO NOTHING RETURNING id::text`, login, passwordHash).Scan(&staffID)
	if errors.Is(err, pgx.ErrNoRows) {
		// A repeat run cannot silently rotate the password or change assigned roles.
		return nil
	}
	if err != nil {
		return err
	}
	for _, role := range roles {
		if _, err := tx.Exec(ctx, `INSERT INTO staff_roles(staff_id, role_code) VALUES ($1, $2)`, staffID, role); err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
