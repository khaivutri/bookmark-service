package jwtutils

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JWTValidator defines the interface for validating JSON Web Tokens.
type JWTValidator interface {
	// ValidateJWT verifies the token string signature using RSA public key and extracts claims.
	ValidateJWT(tokenString string) (jwt.MapClaims, error)
}

type jwtValidator struct {
	publicKey *rsa.PublicKey
}

// NewJWTValidator constructs a new JWT validator using an RSA public key.
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
// ValidateJWT verifies the token string signature using RSA public key and extracts claims.
func (j *jwtValidator) ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, ErrInvalidToken
		}
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