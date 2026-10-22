package auth

import (
	"context"
	"errors"
	"testing"

	"ayurflow/backend/internal/config"
)

func TestLoginStopsWaitingWhenRequestIsCanceled(t *testing.T) {
	service := NewService(nil, config.AuthPolicy{})
	for range cap(service.verifySlots) {
		service.verifySlots <- struct{}{}
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, _, err := service.Login(ctx, "synthetic", "password", ""); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled login = %v", err)
	}
}

func TestSessionSecretIsOpaqueAndHashed(t *testing.T) {
	first, firstHash, err := newSecret()
	if err != nil {
		t.Fatal(err)
	}
	second, secondHash, err := newSecret()
	if err != nil {
		t.Fatal(err)
	}
	if first == second || firstHash == secondHash {
		t.Fatal("session secrets must be independently random")
	}
	if first == string(firstHash[:]) {
		t.Fatal("raw token must differ from stored hash")
	}
	if actual, ok := secretHash(first); !ok || actual != firstHash {
		t.Fatal("valid token did not match its stored hash")
	}
	for _, input := range []string{"", "invalid", first[:len(first)-1]} {
		if _, ok := secretHash(input); ok {
			t.Fatalf("accepted malformed token %q", input)
		}
	}
}
