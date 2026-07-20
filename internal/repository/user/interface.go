package user

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) 
}
