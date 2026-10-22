//go:build integration

package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"os"
	"testing"
	"time"

	"ayurflow/backend/db"
	"ayurflow/backend/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AUTH_TEST_DATABASE_URL must name a disposable PostgreSQL database.
func TestSessionLifecycleIntegration(t *testing.T) {
	url := os.Getenv("AUTH_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("AUTH_TEST_DATABASE_URL is not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	if err := db.Migrate(ctx, pool); err != nil {
		t.Fatal(err)
	}
	password := "synthetic-password-123"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatal(err)
	}
	login := "auth-test-" + time.Now().UTC().Format("20060102150405.000000000")
	var staffID string
	err = pool.QueryRow(ctx, `INSERT INTO staff_accounts(login_identifier, password_hash)
		VALUES ($1, $2) RETURNING id::text`, login, hash).Scan(&staffID)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		cleanup, stop := context.WithTimeout(context.Background(), 10*time.Second)
		defer stop()
		key := sha256.Sum256([]byte(login))
		pool.Exec(cleanup, `DELETE FROM staff_sessions WHERE staff_id=$1`, staffID)
		pool.Exec(cleanup, `DELETE FROM auth_login_attempts WHERE identifier_hash=$1`, key[:])
		pool.Exec(cleanup, `DELETE FROM staff_accounts WHERE id=$1`, staffID)
	}()
	policy := config.AuthPolicy{
		Enabled: true, IdleTimeout: 15 * time.Minute, AbsoluteTimeout: time.Hour,
		AttemptLimit: 2, AttemptWindow: 10 * time.Minute, BlockDuration: 5 * time.Minute,
	}
	service := NewService(pool, policy)
	base := time.Now().UTC()
	service.clock = func() time.Time { return base }
	if _, _, err := service.Login(ctx, login, "wrong", ""); !errors.Is(err, ErrCredentials) {
		t.Fatalf("wrong password: %v", err)
	}
	if _, _, err := service.Login(ctx, login, "wrong", ""); !errors.Is(err, ErrCredentials) {
		t.Fatalf("second wrong password: %v", err)
	}
	if _, _, err := service.Login(ctx, login, password, ""); !errors.Is(err, ErrCredentials) {
		t.Fatalf("blocked login: %v", err)
	}
	service.clock = func() time.Time { return base.Add(6 * time.Minute) }
	_, token, err := service.Login(ctx, login, password, "")
	if err != nil {
		t.Fatal(err)
	}
	_, replacement, err := service.Login(ctx, login, password, token)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.Validate(ctx, token); !errors.Is(err, ErrSession) {
		t.Fatalf("replaced token: %v", err)
	}
	token = replacement
	session, err := service.Validate(ctx, token)
	if err != nil || session.CSRFToken == "" {
		t.Fatalf("valid session: %v", err)
	}
	if err := service.Logout(ctx, token, "wrong"); !errors.Is(err, ErrSession) {
		t.Fatalf("wrong CSRF: %v", err)
	}
	if err := service.Logout(ctx, token, session.CSRFToken); err != nil {
		t.Fatalf("logout: %v", err)
	}
	if _, err := service.Validate(ctx, token); !errors.Is(err, ErrSession) {
		t.Fatalf("logged out token: %v", err)
	}
	_, token, err = service.Login(ctx, login, password, "")
	if err != nil {
		t.Fatal(err)
	}
	service.clock = func() time.Time { return base.Add(22 * time.Minute) }
	if _, err := service.Validate(ctx, token); !errors.Is(err, ErrSession) {
		t.Fatalf("idle expired token: %v", err)
	}
	service.clock = func() time.Time { return base.Add(23 * time.Minute) }
	_, token, err = service.Login(ctx, login, password, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pool.Exec(ctx, `UPDATE staff_accounts SET active=false WHERE id=$1`, staffID); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Validate(ctx, token); !errors.Is(err, ErrSession) {
		t.Fatalf("disabled account token: %v", err)
	}
	if _, err := pool.Exec(ctx, `UPDATE staff_accounts SET active=true WHERE id=$1`, staffID); err != nil {
		t.Fatal(err)
	}
	_, token, err = service.Login(ctx, login, password, "")
	if err != nil {
		t.Fatal(err)
	}
	service.clock = func() time.Time { return base.Add(84 * time.Minute) }
	if _, err := service.Validate(ctx, token); !errors.Is(err, ErrSession) {
		t.Fatalf("absolute expired token: %v", err)
	}
}
