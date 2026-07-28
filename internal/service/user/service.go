package user

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository/user"
	"github.com/khaivutri/bookmark-service/pkg/jwtutils"
	"github.com/khaivutri/bookmark-service/pkg/utils"
)

type Service interface {
	CreateUser(ctx context.Context, userName, displayName, password, email string) (*model.User, error)
	UpdateSelfInfo(ctx context.Context, userID, displayName, email string) (error)
	
	Login(ctx context.Context, userName, password string) (string, error)

	GetSelfInfo(ctx context.Context, userID string) (*model.User, error)
}

type service struct {
	repo 				user.Repository
	hasher	 			utils.Hasher
	jwtGenerator 		jwtutils.JWTGenerator
}

func NewService(repo user.Repository, hasher utils.Hasher, jwtGenerator jwtutils.JWTGenerator) Service {
	return &service{	repo: 	repo, 
						hasher: hasher,
						jwtGenerator: jwtGenerator,
					}
}