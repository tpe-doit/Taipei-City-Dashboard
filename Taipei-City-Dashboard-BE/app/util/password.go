package util

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	passwordMemory     = 19 * 1024
	passwordTime       = 2
	passwordSaltLen    = 16
	passwordKeyLen     = 32
	passwordHashPrefix = "$argon2id$v=19$m=19456,t=2,p=1$"
)

// HashPassword uses a random salt and a memory-hard password hash. HashString
// remains separate because SSO uses it for deterministic identity lookups.
func HashPassword(password string) (string, error) {
	salt := make([]byte, passwordSaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, passwordTime, passwordMemory, 1, passwordKeyLen)
	return passwordHashPrefix + base64.RawStdEncoding.EncodeToString(salt) + "$" + base64.RawStdEncoding.EncodeToString(key), nil
}

// VerifyPassword accepts the current format and, only for existing accounts,
// the old SHA-256 format. A successful legacy check must be upgraded before
// completing login. New passwords must always use HashPassword.
func VerifyPassword(password, encoded string) (valid, needsUpgrade bool) {
	if strings.HasPrefix(encoded, passwordHashPrefix) {
		// Only the fixed, bounded parameters emitted by HashPassword are accepted.
		if len(encoded) != len(passwordHashPrefix)+22+1+43 {
			return false, false
		}
		parts := strings.Split(strings.TrimPrefix(encoded, passwordHashPrefix), "$")
		if len(parts) != 2 || len(parts[0]) != 22 || len(parts[1]) != 43 {
			return false, false
		}
		salt, err := base64.RawStdEncoding.Strict().DecodeString(parts[0])
		if err != nil || len(salt) != passwordSaltLen {
			return false, false
		}
		stored, err := base64.RawStdEncoding.Strict().DecodeString(parts[1])
		if err != nil || len(stored) != passwordKeyLen {
			return false, false
		}
		key := argon2.IDKey([]byte(password), salt, passwordTime, passwordMemory, 1, passwordKeyLen)
		return subtle.ConstantTimeCompare(key, stored) == 1, false
	}

	if len(encoded) != sha256.Size*2 {
		return false, false
	}
	stored, err := hex.DecodeString(encoded)
	if err != nil || encoded != strings.ToLower(encoded) {
		return false, false
	}
	// Compatibility only: never store this digest for a new password.
	digest := sha256.Sum256([]byte(password))
	valid = subtle.ConstantTimeCompare(digest[:], stored) == 1
	return valid, valid
}
