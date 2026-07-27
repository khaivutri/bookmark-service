package jwtutils

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type JWTValidator interface {
	ValidateJWT(tokenString string) (jwt.MapClaims, error)
}

type jwtValidator struct {
	publicKey *rsa.PublicKey
}

func NewJWTValidator(publicKeyPath string) (JWTValidator, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtValidator{
		publicKey: publicKey,
	}, nil
}

var (
	ErrInvalidToken = errors.New("invalid token")	
	ErrExtractToken = errors.New("failed to extract token")
)
func (j *jwtValidator) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return j.publicKey, nil
	})

	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if payload, ok := token.Claims.(jwt.MapClaims); ok {
		return payload, nil
	}

	return nil, ErrExtractToken
}