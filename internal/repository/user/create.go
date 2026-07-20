package user

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
)

func (r *sqlRepository) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		return nil, err
	}	
	return newUser, nil
}