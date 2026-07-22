package user

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository/user"
	"github.com/khaivutri/bookmark-service/pkg/utils"
)

type Service interface {
	CreateUser(ctx context.Context, userName, displayName, password, email string) (*model.User, error)
}

type service struct {
	repo user.Repository
	hasher utils.Hasher
}

func NewService(repo user.Repository, hasher utils.Hasher) Service {
	return &service{	repo: 	repo, 
						hasher: hasher,
					}
}