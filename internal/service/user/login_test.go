package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/internal/model"
	mockRepo "github.com/khaivutri/bookmark-service/internal/repository/user/mocks"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	mockJWT "github.com/khaivutri/bookmark-service/pkg/jwtutils/mocks"
	mockHasher "github.com/khaivutri/bookmark-service/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	stubUser := &model.User{
		ID:       "69c11072-9af6-4a5f-81d5-239d96154d5e",
		UserName: "johndoe",
		Password: "hashed-password",
		Email:    "john.doe@example.com",
	}
	repoErr := errors.New("database connection lost")
	jwtErr := errors.New("sign token failed")

	testCases := []struct {
		name string

		setupMockRepo         func(ctx context.Context, t *testing.T) *mockRepo.Repository
		setupMockHasher       func(t *testing.T) *mockHasher.Hasher
		setupMockJWTGenerator func(t *testing.T) *mockJWT.JWTGenerator

		expectedToken string
		expectedError error
	}{
		{
			name: "success - login returns generated jwt token",

			setupMockRepo: func(ctx context.Context, t *testing.T) *mockRepo.Repository {
				repo := mockRepo.NewRepository(t)
				repo.On("GetUserByUserName", ctx, "johndoe").Return(stubUser, nil).Once()
				return repo
			},
			setupMockHasher: func(t *testing.T) *mockHasher.Hasher {
				hasher := mockHasher.NewHasher(t)
				hasher.On("Compare", "hashed-password", "Password123@").Return(true).Once()
				return hasher
			},
			setupMockJWTGenerator: func(t *testing.T) *mockJWT.JWTGenerator {
				jwtGenerator := mockJWT.NewJWTGenerator(t)
				jwtGenerator.On("GenerateJWT", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					issuedAt, ok := claims["iat"].(int64)
					if !ok {
						return false
					}

					expiresAt, ok := claims["exp"].(int64)
					if !ok {
						return false
					}

					now := time.Now().Unix()
					return claims["sub"] == stubUser.ID &&
						claims["email"] == stubUser.Email &&
						issuedAt >= now-5 &&
						issuedAt <= now+5 &&
						expiresAt-issuedAt == int64(tokenDuration.Seconds())
				})).Return("jwt-token-string", nil).Once()
				return jwtGenerator
			},

			expectedToken: "jwt-token-string",
		},
		{
			name: "error - username does not exist",

			setupMockRepo: func(ctx context.Context, t *testing.T) *mockRepo.Repository {
				repo := mockRepo.NewRepository(t)
				repo.On("GetUserByUserName", ctx, "johndoe").Return(nil, dbutils.ErrRecordNotFound).Once()
				return repo
			},
			setupMockHasher: func(t *testing.T) *mockHasher.Hasher {
				return mockHasher.NewHasher(t)
			},
			setupMockJWTGenerator: func(t *testing.T) *mockJWT.JWTGenerator {
				return mockJWT.NewJWTGenerator(t)
			},

			expectedError: ErrInvalidCredential,
		},
		{
			name: "error - repository returns unexpected error",

			setupMockRepo: func(ctx context.Context, t *testing.T) *mockRepo.Repository {
				repo := mockRepo.NewRepository(t)
				repo.On("GetUserByUserName", ctx, "johndoe").Return(nil, repoErr).Once()
				return repo
			},
			setupMockHasher: func(t *testing.T) *mockHasher.Hasher {
				return mockHasher.NewHasher(t)
			},
			setupMockJWTGenerator: func(t *testing.T) *mockJWT.JWTGenerator {
				return mockJWT.NewJWTGenerator(t)
			},

			expectedError: repoErr,
		},
		{
			name: "error - password does not match",

			setupMockRepo: func(ctx context.Context, t *testing.T) *mockRepo.Repository {
				repo := mockRepo.NewRepository(t)
				repo.On("GetUserByUserName", ctx, "johndoe").Return(stubUser, nil).Once()
				return repo
			},
			setupMockHasher: func(t *testing.T) *mockHasher.Hasher {
				hasher := mockHasher.NewHasher(t)
				hasher.On("Compare", "hashed-password", "Password123@").Return(false).Once()
				return hasher
			},
			setupMockJWTGenerator: func(t *testing.T) *mockJWT.JWTGenerator {
				return mockJWT.NewJWTGenerator(t)
			},

			expectedError: ErrInvalidCredential,
		},
		{
			name: "error - jwt generator fails",

			setupMockRepo: func(ctx context.Context, t *testing.T) *mockRepo.Repository {
				repo := mockRepo.NewRepository(t)
				repo.On("GetUserByUserName", ctx, "johndoe").Return(stubUser, nil).Once()
				return repo
			},
			setupMockHasher: func(t *testing.T) *mockHasher.Hasher {
				hasher := mockHasher.NewHasher(t)
				hasher.On("Compare", "hashed-password", "Password123@").Return(true).Once()
				return hasher
			},
			setupMockJWTGenerator: func(t *testing.T) *mockJWT.JWTGenerator {
				jwtGenerator := mockJWT.NewJWTGenerator(t)
				jwtGenerator.On("GenerateJWT", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					return claims["sub"] == stubUser.ID && claims["email"] == stubUser.Email
				})).Return("", jwtErr).Once()
				return jwtGenerator
			},

			expectedError: jwtErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repo := tc.setupMockRepo(ctx, t)
			hasher := tc.setupMockHasher(t)
			jwtGenerator := tc.setupMockJWTGenerator(t)

			svc := NewService(repo, hasher, jwtGenerator)

			token, err := svc.Login(ctx, "johndoe", "Password123@")

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Empty(t, token)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedToken, token)
		})
	}
}
