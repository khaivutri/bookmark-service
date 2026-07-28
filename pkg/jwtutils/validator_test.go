package jwtutils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJWTValidator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		path       func(t *testing.T) string
		wantErr    bool
		wantPublic bool
	}{
		{
			name: "loads an RSA public key",
			path: func(t *testing.T) string {
				privateKey := generateRSAKey(t)
				return writePublicKey(t, &privateKey.PublicKey)
			},
			wantPublic: true,
		},
		{
			name: "returns an error when the file does not exist",
			path: func(t *testing.T) string {
				return t.TempDir() + string(os.PathSeparator) + "missing.pem"
			},
			wantErr: true,
		},
		{
			name: "returns an error for invalid PEM data",
			path: func(t *testing.T) string {
				path := t.TempDir() + string(os.PathSeparator) + "invalid.pem"
				require.NoError(t, os.WriteFile(path, []byte("not a public key"), 0o600))
				return path
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			validator, err := NewJWTValidator(tc.path(t))

			if tc.wantErr {
				require.Error(t, err)
				assert.Nil(t, validator)
				return
			}

			require.NoError(t, err)
			if tc.wantPublic {
				assert.NotNil(t, validator.(*jwtValidator).publicKey)
			}
		})
	}
}

func TestJWTValidator_ValidateJWT(t *testing.T) {
	t.Parallel()

	privateKey := generateRSAKey(t)
	validator := &jwtValidator{publicKey: &privateKey.PublicKey}

	validToken := signToken(t, jwt.SigningMethodRS256, privateKey, jwt.MapClaims{
		"sub":  "user-123",
		"role": "admin",
	})

	otherKey := generateRSAKey(t)
	wrongSignatureToken := signToken(t, jwt.SigningMethodRS256, otherKey, jwt.MapClaims{
		"sub": "user-123",
	})
	hmacToken := signToken(t, jwt.SigningMethodHS256, []byte("shared-secret"), jwt.MapClaims{
		"sub": "user-123",
	})

	testCases := []struct {
		name       string
		token      string
		wantClaims jwt.MapClaims
		wantErr    error
	}{
		{
			name:       "returns claims for a valid token",
			token:      validToken,
			wantClaims: jwt.MapClaims{"sub": "user-123", "role": "admin"},
		},
		{
			name:    "rejects a malformed token",
			token:   "not-a-jwt",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "rejects a token with an invalid signature",
			token:   wrongSignatureToken,
			wantErr: ErrInvalidToken,
		},
		{
			name:    "rejects a non-RSA signing method",
			token:   hmacToken,
			wantErr: ErrInvalidToken,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			claims, err := validator.ValidateJWT(tc.token)

			if tc.wantErr != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, tc.wantErr))
				assert.Nil(t, claims)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantClaims, claims)
		})
	}
}

func signToken(t *testing.T, method jwt.SigningMethod, key any, claims jwt.MapClaims) string {
	t.Helper()

	token := jwt.NewWithClaims(method, claims)
	tokenString, err := token.SignedString(key)
	require.NoError(t, err)
	return tokenString
}

func writePublicKey(t *testing.T, key *rsa.PublicKey) string {
	t.Helper()

	data := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(key),
	})
	path := t.TempDir() + string(os.PathSeparator) + "public.pem"
	require.NoError(t, os.WriteFile(path, data, 0o600))
	return path
}
