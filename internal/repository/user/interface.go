package user

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
)

// Repository defines the data storage operations for users.
type Repository interface {
	// CreateUser inserts a user record into the database.
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) 
	// UpdateUser saves user changes back to the database.
	UpdateUser(ctx context.Context, user *model.User) (error)

	// GetUserByUserName queries a user record by username.
	GetUserByUserName(ctx context.Context, userName string) (*model.User, error)
	// GetUserByID queries a user record by userID.
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
}
