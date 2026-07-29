package user

import (
	"context"
	"errors"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository/user/mocks"

	"github.com/stretchr/testify/mock"
)

var sampleUser = &model.User{
	ID:          "b649b57b-b7b6-44e4-a233-74147ecf56ee",
	DisplayName: "Test1",
	UserName:    "test1",
	Password:    "pwd123",
	Email:       "test1@example.com",
}

func newTestService(repo *mocks.Repository) *service {
	return &service{repo: repo}
}

type updateSelfInfoTestCase struct {
	name             string
	setupMockRepo    func(t *testing.T) *mocks.Repository
	inputUserID      string
	inputDisplayName string
	inputEmail       string
	expectedErr      error
}

func runUpdateSelfInfoTest(t *testing.T, tc updateSelfInfoTestCase) {
	t.Helper()
	t.Parallel()

	repo := tc.setupMockRepo(t)
	svc := newTestService(repo)

	err := svc.UpdateSelfInfo(context.Background(), tc.inputUserID, tc.inputDisplayName, tc.inputEmail)

	if tc.expectedErr != nil {
		if !errors.Is(err, tc.expectedErr) {
			t.Errorf("expected error %v, got %v", tc.expectedErr, err)
		}
		return
	}

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSelf_GetSelfInfo(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name          string
		setupMockRepo func(t *testing.T) *mocks.Repository
		inputUserID   string

		expectedUser *model.User
		expectedErr  error
	}{
		{
			name: "success - repository returns user",
			setupMockRepo: func(t *testing.T) *mocks.Repository {
				repo := mocks.NewRepository(t)
				repo.On("GetUserByID", mock.Anything, sampleUser.ID).
					Return(sampleUser, nil).
					Once()
				return repo
			},
			inputUserID:  sampleUser.ID,
			expectedUser: sampleUser,
			expectedErr:  nil,
		},
		{
			name: "repository returns error - maps to ErrUserNotFound",
			setupMockRepo: func(t *testing.T) *mocks.Repository {
				repo := mocks.NewRepository(t)
				repo.On("GetUserByID", mock.Anything, "unknown-id").
					Return(nil, errors.New("record not found")).
					Once()
				return repo
			},
			inputUserID:  "unknown-id",
			expectedUser: nil,
			expectedErr:  ErrUserNotFound,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := tc.setupMockRepo(t)
			svc := newTestService(repo)

			user, err := svc.GetSelfInfo(context.Background(), tc.inputUserID)

			if tc.expectedErr != nil {
				if !errors.Is(err, tc.expectedErr) {
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
			if user != tc.expectedUser {
				t.Errorf("expected user %+v, got %+v", tc.expectedUser, user)
			}
		})
	}
}

func TestSelf_UpdateSelfInfo(t *testing.T) {
	t.Parallel()

	testcases := []updateSelfInfoTestCase{
		{
			name: "success - updates display name and email, keeps other fields",
			setupMockRepo: func(t *testing.T) *mocks.Repository {
				repo := mocks.NewRepository(t)

				existing := &model.User{
					ID:          "user-1",
					DisplayName: "Old Name",
					UserName:    "user1",
					Password:    "pwd123",
					Email:       "old@example.com",
				}

				repo.On("GetUserByID", mock.Anything, "user-1").
					Return(existing, nil).
					Once()

		
				repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
					return u.ID == "user-1" &&
						u.DisplayName == "New Name" &&
						u.Email == "new@example.com" &&
						u.UserName == "user1" &&
						u.Password == "pwd123"
				})).Return(nil).Once()

				return repo
			},
			inputUserID:      "user-1",
			inputDisplayName: "New Name",
			inputEmail:       "new@example.com",
			expectedErr:      nil,
		},
		{
			name: "user not found - returns ErrUserNotFound and never calls UpdateUser",
			setupMockRepo: func(t *testing.T) *mocks.Repository {
				repo := mocks.NewRepository(t)
				repo.On("GetUserByID", mock.Anything, "unknown-id").
					Return(nil, errors.New("record not found")).
					Once()
				
				return repo
			},
			inputUserID:      "unknown-id",
			inputDisplayName: "New Name",
			inputEmail:       "new@example.com",
			expectedErr:      ErrUserNotFound,
		},
		{
			name: "repository update fails - returns ErrFailUpdateUser",
			setupMockRepo: func(t *testing.T) *mocks.Repository {
				repo := mocks.NewRepository(t)

				existing := &model.User{
					ID:          "user-3",
					DisplayName: "Old Name",
					UserName:    "user3",
					Email:       "old3@example.com",
				}

				repo.On("GetUserByID", mock.Anything, "user-3").
					Return(existing, nil).
					Once()
				repo.On("UpdateUser", mock.Anything, mock.Anything).
					Return(errors.New("db write error")).
					Once()

				return repo
			},
			inputUserID:      "user-3",
			inputDisplayName: "New Name",
			inputEmail:       "new3@example.com",
			expectedErr:      ErrFailUpdateUser,
		},
		{
			name: "empty display name and email - still forwarded to repository as-is",
			setupMockRepo: func(t *testing.T) *mocks.Repository {
				repo := mocks.NewRepository(t)

				existing := &model.User{
					ID:          "user-4",
					DisplayName: "Old Name",
					UserName:    "user4",
					Email:       "old4@example.com",
				}

				repo.On("GetUserByID", mock.Anything, "user-4").
					Return(existing, nil).
					Once()
				repo.On("UpdateUser", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
					return u.DisplayName == "" && u.Email == ""
				})).Return(nil).Once()

				return repo
			},
			inputUserID:      "user-4",
			inputDisplayName: "",
			inputEmail:       "",
			expectedErr:      nil,
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) { runUpdateSelfInfoTest(t, tc) })
	}
}
