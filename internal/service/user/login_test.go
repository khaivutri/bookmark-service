package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/internal/model"
	mockRepo "github.com/khaivutri/bookmark-service/internal/repository/user/mocks"
	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	mockJWT "github.com/khaivutri/bookmark-service/pkg/jwtutils/mocks"
	mockHasher "github.com/khaivutri/bookmark-service/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	stubUser := &model.User{
		Base:     fixture.GetTestBase("69c11072-9af6-4a5f-81d5-239d96154d5e"),
		UserName: "johndoe",
		Password: "hashed-password",
		Email:    "john.doe@example.com",
	}
	repoErr := errors.New("database connection lost")
	jwtErr := errors.New("sign token failed")

	testCases := []struct {
		name string

		setupMocks func(ctx context.Context, repo *mockRepo.Repository, hasher *mockHasher.Hasher, generator *mockJWT.JWTGenerator)

		expectedToken string
		expectedError error
	}{
		{
			name: "success - login returns generated jwt token",

			setupMocks: func(ctx context.Context, repo *mockRepo.Repository, hasher *mockHasher.Hasher, generator *mockJWT.JWTGenerator) {
				repo.On("GetUserByUserName", ctx, "johndoe").Return(stubUser, nil).Once()
				hasher.On("Compare", "hashed-password", "Password123@").Return(true).Once()
				generator.On("GenerateJWT", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					issuedAt, ok := claims["iat"].(int64)
					if !ok {
						return false
					}

					expiresAt, ok := claims["exp"].(int64)
					if !ok {
						return false
					}

					now := time.Now().Unix()
					return claims["sub"] == stubUser.Base.ID &&
						claims["email"] == stubUser.Email &&
						issuedAt >= now-5 &&
						issuedAt <= now+5 &&
						expiresAt-issuedAt == int64(tokenDuration.Seconds())
				})).Return("jwt-token-string", nil).Once()
			},

			expectedToken: "jwt-token-string",
		},
		{
			name: "error - username does not exist",

			setupMocks: func(ctx context.Context, repo *mockRepo.Repository, _ *mockHasher.Hasher, _ *mockJWT.JWTGenerator) {
				repo.On("GetUserByUserName", ctx, "johndoe").Return(nil, dbutils.ErrRecordNotFound).Once()
			},

			expectedError: ErrInvalidCredential,
		},
		{
			name: "error - repository returns unexpected error",

			setupMocks: func(ctx context.Context, repo *mockRepo.Repository, _ *mockHasher.Hasher, _ *mockJWT.JWTGenerator) {
				repo.On("GetUserByUserName", ctx, "johndoe").Return(nil, repoErr).Once()
			},

			expectedError: repoErr,
		},
		{
			name: "error - password does not match",

			setupMocks: func(ctx context.Context, repo *mockRepo.Repository, hasher *mockHasher.Hasher, _ *mockJWT.JWTGenerator) {
				repo.On("GetUserByUserName", ctx, "johndoe").Return(stubUser, nil).Once()
				hasher.On("Compare", "hashed-password", "Password123@").Return(false).Once()
			},

			expectedError: ErrInvalidCredential,
		},
		{
			name: "error - jwt generator fails",

			setupMocks: func(ctx context.Context, repo *mockRepo.Repository, hasher *mockHasher.Hasher, generator *mockJWT.JWTGenerator) {
				repo.On("GetUserByUserName", ctx, "johndoe").Return(stubUser, nil).Once()
				hasher.On("Compare", "hashed-password", "Password123@").Return(true).Once()
				generator.On("GenerateJWT", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					return claims["sub"] == stubUser.Base.ID && claims["email"] == stubUser.Email
				})).Return("", jwtErr).Once()
			},

			expectedError: jwtErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repo := mockRepo.NewRepository(t)
			hasher := mockHasher.NewHasher(t)
			jwtGenerator := mockJWT.NewJWTGenerator(t)
			tc.setupMocks(ctx, repo, hasher, jwtGenerator)

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
