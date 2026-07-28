package user

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
)

const tokenDuration = 24 * time.Hour 

var ( 
	ErrInvalidCredential = errors.New("invalid credential")
)

func (u *service) Login(ctx context.Context, userName, password string) (string, error) {
	// check username 
	user, err := u.repo.GetUserByUserName(ctx, userName)
	switch {
		case errors.Is(err, dbutils.ErrRecordNotFound):
			return "", ErrInvalidCredential
		case err == nil:
		default:
			return "", err
	}
	// compare hashed password
	if !u.hasher.Compare(user.Password, password) {
		return "", ErrInvalidCredential
	}
	// if match -> generate token, return it
	tokenContent := jwt.MapClaims{
		"sub": 			user.ID,
		"email": 		user.Email,
		"iat": 			time.Now().Unix(),
		"exp": 			time.Now().Add(tokenDuration).Unix(),
	}
	
	tokenString, err := u.jwtGenerator.GenerateJWT(tokenContent)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}