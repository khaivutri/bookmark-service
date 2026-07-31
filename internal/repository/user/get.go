package user

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
)

// GetUserByUserName queries a user record by username.
func (r *sqlRepository) GetUserByUserName(ctx context.Context, userName string) (*model.User, error) {
	user := &model.User{}
	
	err := r.db.WithContext(ctx).Where("username = ?", userName).First(user).Error
	if err != nil {
		return nil, dbutils.ParseDBError(err)
	}
	
	return user, nil
}

// GetUserByID queries a user record by userID.
func (r *sqlRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	user := &model.User{}
	
	err := r.db.WithContext(ctx).Where("id = ?", userID).First(user).Error
	if err != nil {
		return nil, dbutils.ParseDBError(err)
	}
	
	return user, nil
}

// UpdateUser saves user changes back to the database.
func (r *sqlRepository) UpdateUser(ctx context.Context, user *model.User) (error) {
	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		return dbutils.ParseDBError(err)
	}
	
	return nil
}