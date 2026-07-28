package user

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) 
	
	GetUserByUserName(ctx context.Context, userName string) (*model.User, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}
