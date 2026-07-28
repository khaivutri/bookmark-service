package user

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
)

func (r *sqlRepository)GetUserByUserName(ctx context.Context, userName string) (*model.User, error) {
	user := &model.User{}
	
	err := r.db.WithContext(ctx).Where("user_name = ?", userName).First(user).Error
	if err != nil {
		return nil, dbutils.ParseDBError(err)
	}
	
	return user, nil
}