package user

import (
	"context"
	"errors"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	
	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"

	"gorm.io/gorm"
)
var (
	existingUser1 = &model.User{
		ID:          "b649b57b-b7b6-44e4-a233-74147ecf56ee",
		DisplayName: "Test1",
		UserName:    "test1",
		Password:    "pwd123",
		Email:       "test1@example.com",
	}
	existingUser2 = &model.User{
		ID:          "b649b57b-b7b6-44e4-a233-74147ecf56ef",
		DisplayName: "Test2",
		UserName:    "test2",
		Password:    "pwd123",
		Email:       "test2@example.com",
	}
)


func newTestRepository(t *testing.T) (*sqlRepository, *gorm.DB) {
	t.Helper()

	db := fixture.NewFixture(t, &fixture.UserCommonTest{})

	return &sqlRepository{db: db}, db
}

func assertUserEqual(t *testing.T, want, got *model.User) {
	t.Helper()

	if got == nil {
		t.Fatalf("expected user %+v, got nil", want)
	}

	if got.ID != want.ID {
		t.Errorf("ID mismatch: want %q, got %q", want.ID, got.ID)
	}
	if got.UserName != want.UserName {
		t.Errorf("UserName mismatch: want %q, got %q", want.UserName, got.UserName)
	}
	if got.DisplayName != want.DisplayName {
		t.Errorf("DisplayName mismatch: want %q, got %q", want.DisplayName, got.DisplayName)
	}
	if got.Email != want.Email {
		t.Errorf("Email mismatch: want %q, got %q", want.Email, got.Email)
	}
}

func TestSqlRepository_GetUserByUserName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		inputUserName string
		expectedUser  *model.User
		expectedErr   error
	}{
		{
			name:          "existing user - returns user successfully",
			inputUserName: "test1",
			expectedUser:  existingUser1,
			expectedErr:   nil,
		},
		{
			name:          "another existing user - returns user successfully",
			inputUserName: "test2",
			expectedUser:  existingUser2,
			expectedErr:   nil,
		},
		{
			name:          "user not found - returns parsed not found error",
			inputUserName: "non_existent_user",
			expectedUser:  nil,
			expectedErr:   dbutils.ParseDBError(gorm.ErrRecordNotFound),
		},
		{
			name:          "empty user name - returns not found error",
			inputUserName: "",
			expectedUser:  nil,
			expectedErr:   dbutils.ParseDBError(gorm.ErrRecordNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, _ := newTestRepository(t)

			user, err := repo.GetUserByUserName(context.Background(), tc.inputUserName)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.expectedErr)
				}
				if !errors.Is(err, tc.expectedErr) && err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tc.expectedErr, err)
				}
				if user != nil {
					t.Errorf("expected nil user on error, got %+v", user)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertUserEqual(t, tc.expectedUser, user)
		})
	}
}

func TestSqlRepository_GetUserByID(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		inputUserID  string
		expectedUser *model.User
		expectedErr  error
	}{
		{
			name:         "existing id - returns user successfully",
			inputUserID:  existingUser1.ID,
			expectedUser: existingUser1,
			expectedErr:  nil,
		},
		{
			name:         "another existing id - returns user successfully",
			inputUserID:  existingUser2.ID,
			expectedUser: existingUser2,
			expectedErr:  nil,
		},
		{
			name:         "id not found - returns parsed not found error",
			inputUserID:  "00000000-0000-0000-0000-000000000000",
			expectedUser: nil,
			expectedErr:  dbutils.ParseDBError(gorm.ErrRecordNotFound),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, _ := newTestRepository(t)

			user, err := repo.GetUserByID(t.Context(), tc.inputUserID)

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tc.expectedErr)
				}
				if !errors.Is(err, tc.expectedErr) && err.Error() != tc.expectedErr.Error() {
					t.Errorf("expected error %v, got %v", tc.expectedErr, err)
				}
				if user != nil {
					t.Errorf("expected nil user on error, got %+v", user)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			assertUserEqual(t, tc.expectedUser, user)
		})
	}
}

func TestSqlRepository_UpdateUser(t *testing.T) {
	t.Parallel()

	t.Run("update existing user - persists changes", func(t *testing.T) {
		t.Parallel()

		repo, db := newTestRepository(t)


		updated := *existingUser1
		updated.DisplayName = "Test1 Updated"
		updated.Email = "test1_updated@example.com"

		err := repo.UpdateUser(context.Background(), &updated)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		
		var got model.User
		if err := db.Where("id = ?", existingUser1.ID).First(&got).Error; err != nil {
			t.Fatalf("failed to fetch updated user: %v", err)
		}

		if got.DisplayName != "Test1 Updated" {
			t.Errorf("DisplayName not updated: want %q, got %q", "Test1 Updated", got.DisplayName)
		}
		if got.Email != "test1_updated@example.com" {
			t.Errorf("Email not updated: want %q, got %q", "test1_updated@example.com", got.Email)
		}
		if got.UserName != existingUser1.UserName {
			t.Errorf("UserName should remain unchanged: want %q, got %q", existingUser1.UserName, got.UserName)
		}
	})

	t.Run("update user with new id - creates new record (gorm Save upsert behavior)", func(t *testing.T) {
		t.Parallel()

		repo, db := newTestRepository(t)

		newUser := &model.User{
			ID:          "c1111111-1111-1111-1111-111111111111",
			DisplayName: "New User",
			UserName:    "newuser",
			Password:    "pwd123",
			Email:       "newuser@example.com",
		}

		err := repo.UpdateUser(t.Context(), newUser)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		var got model.User
		if err := db.Where("id = ?", newUser.ID).First(&got).Error; err != nil {
			t.Fatalf("expected new record to be created by Save(), got error: %v", err)
		}
		assertUserEqual(t, newUser, &got)
	})
}