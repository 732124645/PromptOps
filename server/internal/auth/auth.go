// Package auth provides password hashing and session-token helpers.
package auth

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"
)

const (
	iterations = 100_000
	keyLength  = 32
	saltLength = 16
)

// HashPassword returns a "salt:hash" string (both hex-encoded) for storage.
func HashPassword(password string) (string, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	dk, err := pbkdf2.Key(sha256.New, password, salt, iterations, keyLength)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(dk), nil
}

// VerifyPassword reports whether password matches a stored "salt:hash" string.
func VerifyPassword(stored, password string) bool {
	salt, want, err := splitHash(stored)
	if err != nil {
		return false
	}
	got, err := pbkdf2.Key(sha256.New, password, salt, iterations, keyLength)
	if err != nil {
		return false
	}
	return subtle.ConstantTimeCompare(got, want) == 1
}

func splitHash(stored string) (salt, hash []byte, err error) {
	parts := strings.SplitN(stored, ":", 2)
	if len(parts) != 2 {
		return nil, nil, errors.New("malformed password hash")
	}
	if salt, err = hex.DecodeString(parts[0]); err != nil {
		return nil, nil, err
	}
	if hash, err = hex.DecodeString(parts[1]); err != nil {
		return nil, nil, err
	}
	return salt, hash, nil
}

// NewToken returns a random, URL-safe session token.
func NewToken() string {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
