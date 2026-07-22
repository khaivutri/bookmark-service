package user_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	"github.com/khaivutri/bookmark-service/internal/model"
	repoUser "github.com/khaivutri/bookmark-service/internal/repository/user"
	serviceUser "github.com/khaivutri/bookmark-service/internal/service/user"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	utilMocks "github.com/khaivutri/bookmark-service/pkg/utils/mocks"
)

func TestService_CreateUser_IntegrationWithFixture(t *testing.T) {
	ctx := context.Background()

	db := fixture.NewFixture(t, &fixture.UserCommonTest{})
	realRepo := repoUser.NewSqlRepository(db)

	tests := []struct {
		name          			string
		input         			struct {
			userName    	string
			displayName 	string
			password    	string
			email       	string
		}
		setupMock     func(hasher *utilMocks.Hasher)
		expectedUser  *model.User
		expectedError error
	}{
		{
			name: "SUCCESS - Register new user successfully",
			input: struct {
				userName    string
				displayName string
				password    string
				email       string
			}{
				userName:    "newuser123",
				displayName: "New User",
				password:    "password123",
				email:       "newuser@example.com",
			},
			setupMock: func(hasher *utilMocks.Hasher) {
				hasher.On("Hash", "password123").
					Return("hashed_pwd_123", nil).
					Once()
			},
			expectedUser: &model.User{
				UserName:    "newuser123",
				DisplayName: "New User",
				Password:    "hashed_pwd_123",
				Email:       "newuser@example.com",
			},
			expectedError: nil,
		},
		{
			name: "ERROR - Duplicate UserName (Conflict with fixture data 'test1')",
			input: struct {
				userName    string
				displayName string
				password    string
				email       string
			}{
				userName:    "test1", 
				displayName: "Test 1 Duplicate",
				password:    "password123",
				email:       "unique_email@example.com",
			},
			setupMock: func(hasher *utilMocks.Hasher) {
				hasher.On("Hash", "password123").
					Return("hashed_pwd_123", nil).
					Once()
			},
			expectedUser:  nil,
			expectedError: dbutils.ErrDuplicateUserName,
		},
		{
			name: "ERROR - Duplicate Email (Conflict with fixture data 'test2@example.com')",
			input: struct {
				userName    string
				displayName string
				password    string
				email       string
			}{
				userName:    "another_unique_user",
				displayName: "Duplicate Email User",
				password:    "password123",
				email:       "test2@example.com", 
			},
			setupMock: func(hasher *utilMocks.Hasher) {
				hasher.On("Hash", "password123").
					Return("hashed_pwd_123", nil).
					Once()
			},
			expectedUser:  nil,
			expectedError: dbutils.ErrDuplicateEmail,
		},
		{
			name: "ERROR - Password Hashing Failed",
			input: struct {
				userName    string
				displayName string
				password    string
				email       string
			}{
				userName:    "fail_hash_user",
				displayName: "Fail Hash",
				password:    "password123",
				email:       "failhash@example.com",
			},
			setupMock: func(hasher *utilMocks.Hasher) {
				hasher.On("Hash", "password123").
					Return("", errors.New("hash algorithm failed")).
					Once()
			},
			expectedUser:  nil,
			expectedError: errors.New("hash algorithm failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockHasher := utilMocks.NewHasher(t)
			tt.setupMock(mockHasher)

			svc := serviceUser.NewService(realRepo, mockHasher)
			createdUser, err := svc.CreateUser(
				ctx,
				tt.input.userName,
				tt.input.displayName,
				tt.input.password,
				tt.input.email,
			)

			// Assert
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Nil(t, createdUser)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, createdUser)
				assert.NotEmpty(t, createdUser.ID)
				assert.Equal(t, tt.expectedUser.UserName, createdUser.UserName)
				assert.Equal(t, tt.expectedUser.DisplayName, createdUser.DisplayName)
				assert.Equal(t, tt.expectedUser.Password, createdUser.Password)
				assert.Equal(t, tt.expectedUser.Email, createdUser.Email)
			}
		})
	}
}