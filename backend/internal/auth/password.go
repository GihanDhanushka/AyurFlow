package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	argonMemoryKiB   uint32 = 64 * 1024
	argonPasses      uint32 = 3
	argonThreads     uint8  = 4
	saltLength              = 16
	keyLength               = 32
	maxPasswordBytes        = 1024
)

var (
	ErrInvalidPassword = errors.New("password must be nonempty and at most 1024 bytes")
	ErrInvalidHash     = errors.New("invalid password hash")
)

// HashPassword returns a PHC-formatted Argon2id hash with a new random salt.
// Password policy belongs to the account workflow; this function only bounds input size.
func HashPassword(password string) (string, error) {
	if len(password) == 0 || len(password) > maxPasswordBytes {
		return "", ErrInvalidPassword
	}
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, argonPasses, argonMemoryKiB, argonThreads, keyLength)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemoryKiB, argonPasses, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt), base64.RawStdEncoding.EncodeToString(key)), nil
}

// VerifyPassword accepts only bounded Argon2id parameters before doing memory-hard work.
// Callers must return generic authentication failures regardless of this result.
func VerifyPassword(password, encoded string) (bool, error) {
	if len(password) == 0 || len(password) > maxPasswordBytes {
		return false, ErrInvalidPassword
	}
	if len(encoded) > 512 {
		return false, ErrInvalidHash
	}
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[0] != "" || parts[1] != "argon2id" || parts[2] != "v=19" {
		return false, ErrInvalidHash
	}
	parameters := strings.Split(parts[3], ",")
	if len(parameters) != 3 {
		return false, ErrInvalidHash
	}
	memory, err := parseParameter(parameters[0], "m=", 19*1024, 256*1024)
	if err != nil {
		return false, ErrInvalidHash
	}
	passes, err := parseParameter(parameters[1], "t=", 2, 10)
	if err != nil {
		return false, ErrInvalidHash
	}
	threads, err := parseParameter(parameters[2], "p=", 1, 8)
	if err != nil || memory < 8*threads {
		return false, ErrInvalidHash
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) < 16 || len(salt) > 64 {
		return false, ErrInvalidHash
	}
	key, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(key) < 16 || len(key) > 64 {
		return false, ErrInvalidHash
	}
	actual := argon2.IDKey([]byte(password), salt, uint32(passes), uint32(memory), uint8(threads), uint32(len(key)))
	return subtle.ConstantTimeCompare(actual, key) == 1, nil
}

func parseParameter(value, prefix string, min, max uint64) (uint64, error) {
	if !strings.HasPrefix(value, prefix) {
		return 0, ErrInvalidHash
	}
	number, err := strconv.ParseUint(strings.TrimPrefix(value, prefix), 10, 64)
	if err != nil || number < min || number > max {
		return 0, ErrInvalidHash
	}
	return number, nil
}
