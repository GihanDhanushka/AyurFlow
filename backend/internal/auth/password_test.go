package auth

import (
	"errors"
	"strings"
	"testing"
)

func TestPasswordHashAndVerify(t *testing.T) {
	first, err := HashPassword("a long synthetic passphrase")
	if err != nil {
		t.Fatal(err)
	}
	second, err := HashPassword("a long synthetic passphrase")
	if err != nil {
		t.Fatal(err)
	}
	if first == second || !strings.HasPrefix(first, "$argon2id$v=19$") {
		t.Fatal("hashes need unique salts and Argon2id encoding")
	}
	for _, tc := range []struct {
		password string
		want     bool
	}{
		{"a long synthetic passphrase", true},
		{"a wrong synthetic passphrase", false},
	} {
		ok, err := VerifyPassword(tc.password, first)
		if err != nil || ok != tc.want {
			t.Fatalf("VerifyPassword: got %v, %v; want %v", ok, err, tc.want)
		}
	}
}

func TestPasswordInputsAreBounded(t *testing.T) {
	for _, value := range []string{"", strings.Repeat("a", 1025)} {
		if _, err := HashPassword(value); !errors.Is(err, ErrInvalidPassword) {
			t.Fatalf("HashPassword: got %v, want ErrInvalidPassword", err)
		}
	}
	for _, encoded := range []string{
		"plaintext",
		strings.Repeat("x", 513),
		"$argon2i$v=19$m=65536,t=3,p=4$c2FsdC0xMjM0NTY3ODkwMTIzNA$a2V5",
		"$argon2id$v=19$m=999999999,t=3,p=4$c2FsdC0xMjM0NTY3ODkwMTIzNA$a2V5",
		"$argon2id$v=19$m=65536,t=0,p=4$c2FsdC0xMjM0NTY3ODkwMTIzNA$a2V5",
	} {
		if _, err := VerifyPassword("a long synthetic passphrase", encoded); !errors.Is(err, ErrInvalidHash) {
			t.Fatalf("VerifyPassword(%q): got %v, want ErrInvalidHash", encoded, err)
		}
	}
}
