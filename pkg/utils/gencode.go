package utils

import (
	"crypto/rand"
	"math/big"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenCode generates secure random alphanumeric codes of a given length.
type GenCode interface {
	// Generate creates a secure random code of the specified length.
	Generate(codeLen int) (string, error)
}

type genCode struct {
}

// NewGenCode constructs a new code generator.
func NewGenCode() GenCode {
	return &genCode{}
}

// Generate generates a secure random code of the specified length.
func (g *genCode) Generate(codeLen int) (string, error) {
	code := make([]byte, codeLen)

	for i := range code {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charSet))))

		if err != nil {
			return "", err
		}

		code[i] = charSet[index.Int64()]
	}
	return string(code), nil
}