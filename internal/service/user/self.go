package user

import (
	"context"
	"errors"

	"github.com/khaivutri/bookmark-service/internal/model"
)
var (
	ErrUserNotFound = errors.New("user not found")
	ErrFailUpdateUser = errors.New("failed to update user")
)

// GetSelfInfo retrieves user information by userID.
func (s *service) GetSelfInfo(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// UpdateSelfInfo updates email and display name info for a user.
func (s *service) UpdateSelfInfo(ctx context.Context, userID, displayName, email string) error {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return ErrUserNotFound
	}
	
	user.DisplayName = displayName
	user.Email = email
	
	err = s.repo.UpdateUser(ctx, user)
	if err != nil {
		return ErrFailUpdateUser
	}
	return nil
}