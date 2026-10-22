package patients

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// These checks implement the provisional synthetic-development rule in DEC02.
// Registration remains unavailable until auth, record access and audit work.
var ErrInvalidPatient = errors.New("patients: invalid patient fields")

type IdentityInput struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type RegistrationInput struct {
	DisplayName string         `json:"displayName"`
	Phone       *string        `json:"phone,omitempty"`
	Identity    *IdentityInput `json:"identity,omitempty"`
}

// ValidateRegistration validates the smallest provisional patient record.
// Clinical review may add fields or revise identity equivalence before AT02.
func ValidateRegistration(input RegistrationInput) error {
	if !utf8.ValidString(input.DisplayName) || strings.TrimSpace(input.DisplayName) == "" {
		return ErrInvalidPatient
	}
	if input.Phone != nil && (!utf8.ValidString(*input.Phone) || strings.TrimSpace(*input.Phone) == "") {
		return ErrInvalidPatient
	}
	if input.Identity != nil {
		if input.Identity.Type != "NIC" && input.Identity.Type != "PASSPORT" {
			return ErrInvalidPatient
		}
		if _, err := CanonicalIdentity(input.Identity.Value); err != nil {
			return err
		}
	}
	return nil
}

// CanonicalIdentity matches patient_identities.canonical_value in migration 000005.
// It trims outer ASCII whitespace, removes ASCII spaces/hyphens and uppercases
// ASCII letters. Other characters remain distinct pending clinic review.
func CanonicalIdentity(value string) (string, error) {
	if !utf8.ValidString(value) {
		return "", ErrInvalidPatient
	}
	value = strings.Trim(value, " \t\n\r")
	value = strings.ReplaceAll(strings.ReplaceAll(value, " ", ""), "-", "")
	var b strings.Builder
	b.Grow(len(value))
	for i := 0; i < len(value); i++ {
		c := value[i]
		if c >= 'a' && c <= 'z' {
			c -= 'a' - 'A'
		}
		b.WriteByte(c)
	}
	if b.Len() == 0 {
		return "", ErrInvalidPatient
	}
	return b.String(), nil
}
