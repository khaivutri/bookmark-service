package user

import (
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqlRepository_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		inputUserName string
		inputEmail    string
		expectedError error
		verifyFunc    func(db *gorm.DB, userName, email string)
	}{
		{
			name: "normal case",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTest{})
			},
			inputUserName: "new_user_name",
			inputEmail:    "new_email@test.com",
			expectedError: nil,
			verifyFunc: func(db *gorm.DB, userName, email string) {
				user := model.User{}
				err := db.First(&user, "user_name = ?", userName).Error

				assert.NoError(t, err)
				assert.Equal(t, userName, user.UserName)
				assert.Equal(t, email, user.Email)
				assert.NotEmpty(t, user.ID)
			},
		},
		{
			name: "duplicate username error",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTest{})
			},
			inputUserName: "test1", // duplicate username in fixture
			inputEmail:    "unique_email@test.com",
			expectedError: dbutils.ErrDuplicateUserName,
			verifyFunc:    nil,
		},
		{
			name: "duplicate email error",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixture.NewFixture(t, &fixture.UserCommonTest{})
			},
			inputUserName: "unique_user_name",
			inputEmail:    "test2@example.com", // 	duplicate email in fixture
			expectedError: dbutils.ErrDuplicateEmail,
			verifyFunc:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			db := tc.setupDB(t)
			repo := NewSqlRepository(db)

			user, err := repo.CreateUser(ctx, &model.User{
				UserName: tc.inputUserName,
				Email:    tc.inputEmail,
			})

			assert.ErrorIs(t, err, tc.expectedError)

			if tc.expectedError == nil {
				assert.NotNil(t, user)
				if tc.verifyFunc != nil {
					tc.verifyFunc(db, tc.inputUserName, tc.inputEmail)
				}
			} else {
				assert.Nil(t, user)
			}
		})
	}
}