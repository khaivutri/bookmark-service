package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)
// Hasher provides methods for hashing passwords and comparing them.
type Hasher interface {
	// Hash encrypts the plain password using bcrypt.
	Hash(password string) (string, error)
	// Compare verifies if a plain password matches a bcrypt hash.
	Compare(hash, password string) bool
}

type hasher struct {
	
}

// NewHasher constructs a new password Hasher.
func NewHasher() Hasher {
	return &hasher{}
}

var ErrHashFailed = errors.New("failed to hash password")
// Hash encrypts the plain password using bcrypt.
func (h *hasher) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashFailed
	}
	return string(hashBytes), nil
}

// Compare verifies if a plain password matches a bcrypt hash.
func (h *hasher) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) 
	return err == nil
}