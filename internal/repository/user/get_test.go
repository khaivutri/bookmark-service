package user

import (
	"context"
	"errors"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/assert"

	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"

	"gorm.io/gorm"
)

var (
	existingUser1 = &model.User{
		Base:        fixture.GetTestBase("b649b57b-b7b6-44e4-a233-74147ecf56ee"),
		DisplayName: "Test1",
		UserName:    "test1",
		Password:    "pwd123",
		Email:       "test1@example.com",
	}
	existingUser2 = &model.User{
		Base:        fixture.GetTestBase("b649b57b-b7b6-44e4-a233-74147ecf56ef"),
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

type userLookupCase struct {
	name         string
	input        string
	expectedUser *model.User
	expectedErr  error
}

type userUpdateCase struct {
	name         string
	inputUser    *model.User
	expectedErr  error
	verifyStored bool
}

func testUserLookup(t *testing.T, lookup func(*sqlRepository, context.Context, string) (*model.User, error), cases []userLookupCase) {
	t.Helper()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, _ := newTestRepository(t)
			got, err := lookup(repo, t.Context(), tc.input)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Nil(t, got, "expected nil user on error")
				return
			}

			assert.NoError(t, err)
			assertUserEqual(t, tc.expectedUser, got)
		})
	}
}

func TestSqlRepository_GetUserByUserName(t *testing.T) {
	t.Parallel()

	testUserLookup(t, func(repo *sqlRepository, ctx context.Context, input string) (*model.User, error) {
		return repo.GetUserByUserName(ctx, input)
	}, []userLookupCase{
		{
			name:         "existing user - returns user successfully",
			input:        "test1",
			expectedUser: existingUser1,
			expectedErr:  nil,
		},
		{
			name:         "another existing user - returns user successfully",
			input:        "test2",
			expectedUser: existingUser2,
			expectedErr:  nil,
		},
		{
			name:         "user not found - returns parsed not found error",
			input:        "non_existent_user",
			expectedUser: nil,
			expectedErr:  dbutils.ErrRecordNotFound,
		},
		{
			name:         "empty user name - returns not found error",
			input:        "",
			expectedUser: nil,
			expectedErr:  dbutils.ErrRecordNotFound,
		},
	})
}

func TestSqlRepository_GetUserByID(t *testing.T) {
	t.Parallel()

	testUserLookup(t, func(repo *sqlRepository, ctx context.Context, input string) (*model.User, error) {
		return repo.GetUserByID(ctx, input)
	}, []userLookupCase{
		{
			name:         "existing id - returns user successfully",
			input:        existingUser1.ID,
			expectedUser: existingUser1,
			expectedErr:  nil,
		},
		{
			name:         "another existing id - returns user successfully",
			input:        existingUser2.ID,
			expectedUser: existingUser2,
			expectedErr:  nil,
		},
		{
			name:         "id not found - returns parsed not found error",
			input:        "00000000-0000-0000-0000-000000000000",
			expectedUser: nil,
			expectedErr:  dbutils.ErrRecordNotFound,
		},
	})
}

func TestSqlRepository_UpdateUser(t *testing.T) {
	t.Parallel()

	testCases := []userUpdateCase{
		{
			name: "update existing user successfully",
			inputUser: &model.User{
				Base:        fixture.GetTestBase(existingUser1.ID),
				DisplayName: "Test1 Updated",
				UserName:    existingUser1.UserName,
				Password:    existingUser1.Password,
				Email:       "test1_updated@example.com",
			},
			verifyStored: true,
		},
		{
			name: "new id creates a user",
			inputUser: &model.User{
				Base:        fixture.GetTestBase("c1111111-1111-1111-1111-111111111111"),
				DisplayName: "New User",
				UserName:    "newuser",
				Password:    "pwd123",
				Email:       "newuser@example.com",
			},
			verifyStored: true,
		},
		{
			name: "duplicate username returns parsed database error",
			inputUser: &model.User{
				Base:        fixture.GetTestBase(existingUser1.ID),
				DisplayName: existingUser1.DisplayName,
				UserName:    existingUser2.UserName,
				Password:    existingUser1.Password,
				Email:       "unique@example.com",
			},
			expectedErr: dbutils.ErrDuplicateUserName,
		},
		{
			name: "duplicate email returns parsed database error",
			inputUser: &model.User{
				Base:        fixture.GetTestBase(existingUser1.ID),
				DisplayName: existingUser1.DisplayName,
				UserName:    "unique_user",
				Password:    existingUser1.Password,
				Email:       existingUser2.Email,
			},
			expectedErr: dbutils.ErrDuplicateEmail,
		},
	}

	testUserUpdates(t, testCases)
}

func testUserUpdates(t *testing.T, testCases []userUpdateCase) {
	t.Helper()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, db := newTestRepository(t)
			err := repo.UpdateUser(t.Context(), tc.inputUser)
			if tc.expectedErr != nil {
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected error %v, got %v", tc.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !tc.verifyStored {
				return
			}

			var got model.User
			if err := db.First(&got, "id = ?", tc.inputUser.Base.ID).Error; err != nil {
				t.Fatalf("failed to fetch stored user: %v", err)
			}
			assertUserEqual(t, tc.inputUser, &got)
		})
	}
}
