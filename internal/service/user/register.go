package user

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
)

func (s *service) CreateUser(ctx context.Context, userName, displayName, password, email string) (*model.User, error) {
	// hash pwd 
	hashPwd, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}
	// create user 
	newUser := &model.User{	
					UserName: 		userName, 
					DisplayName: 	displayName, 
					Password: 		hashPwd, 
					Email: 			email,
				}
	// call repo to create user
	res, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}

	// return user
	return res, nil
}