package user

import (
	"context"
	"errors"

	"github.com/khaivutri/bookmark-service/internal/model"
)
var (
	ErrUserNotFound = errors.New("user not found")
)
func (s *service) GetSelfInfo(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}