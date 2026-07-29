package jwtutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTGenerator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		path        func(t *testing.T) string
		wantErr     bool
		wantPrivate bool
	}{
		{
			name: "loads an RSA private key",
			path: func(t *testing.T) string {
				return writePrivateKey(t)
			},
			wantPrivate: true,
		},
		{
			name: "returns an error when the file does not exist",
			path: func(t *testing.T) string {
				return filepath.Join(t.TempDir(), "missing.pem")
			},
			wantErr: true,
		},
		{
			name: "returns an error for invalid PEM data",
			path: func(t *testing.T) string {
				path := filepath.Join(t.TempDir(), "invalid.pem")
				require.NoError(t, os.WriteFile(path, []byte("not a private key"), 0o600))
				return path
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			generator, err := NewJWTGenerator(tc.path(t))

			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, generator)
				return
			}

			require.NoError(t, err)
			if tc.wantPrivate {
				assert.NotNil(t, generator.(*jwtGenerator).privateKey)
			}
		})
	}
}

func TestJWTGenerator_GenerateJWT(t *testing.T) {
	t.Parallel()

	privateKey := generateRSAKey(t)
	generator := &jwtGenerator{privateKey: privateKey}

	testCases := []struct {
		name   string
		claims jwt.MapClaims
	}{
		{
			name: "generates a token with all claims",
			claims: jwt.MapClaims{
				"sub":  "user-123",
				"role": "admin",
				"exp":  float64(2_000_000_000),
			},
		},
		{
			name:   "generates a token with empty claims",
			claims: jwt.MapClaims{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			tokenString, err := generator.GenerateJWT(tc.claims)
			require.NoError(t, err)
			assert.NotEmpty(t, tokenString)

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
				assert.Equal(t, jwt.SigningMethodRS256, token.Method)
				return &privateKey.PublicKey, nil
			})
			require.NoError(t, err)
			require.True(t, token.Valid)
			assert.Equal(t, tc.claims, token.Claims.(jwt.MapClaims))
		})
	}
}

func generateRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return key
}

func writePrivateKey(t *testing.T) string {
	t.Helper()

	key := generateRSAKey(t)
	data := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	path := filepath.Join(t.TempDir(), "private.pem")
	require.NoError(t, os.WriteFile(path, data, 0o600))
	return path
}
