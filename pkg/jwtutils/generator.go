package jwtutils

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JWTGenerator defines the interface for generating signed JSON Web Tokens.
type JWTGenerator interface {
	// GenerateJWT signs and generates a JWT string using MapClaims and RSA private key.
	GenerateJWT(jwtContent jwt.MapClaims) (string, error)
}

type jwtGenerator struct {
	privateKey *rsa.PrivateKey
}

// NewJWTGenerator constructs a new JWT generator using an RSA private key.
func NewJWTGenerator(privateKeyPath string) (JWTGenerator, error){
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}
	
	return &jwtGenerator{
		privateKey: privateKey,
	}, nil
}

// GenerateJWT signs and generates a JWT string using MapClaims and RSA private key.
func (g *jwtGenerator) GenerateJWT(jwtContent jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtContent)
	
	tokenString, err := token.SignedString(g.privateKey)
	if err != nil {
		return "", err
	}
	
	return tokenString, nil
}