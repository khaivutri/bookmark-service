package utils

import (
	"crypto/rand"
	"math/big"
)

const charSet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type GenCode interface {
	Generate(codeLen int) (string, error)
}

type genCode struct {
}

func NewGenCode() GenCode {
	return &genCode{}
}

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