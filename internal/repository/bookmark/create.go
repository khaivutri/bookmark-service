package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
)

func (r *bookmarkRepo) CreateBookmark(ctx context.Context, newBookmark *model.Bookmark) (*model.Bookmark, error) {
	err := r.db.WithContext(ctx).Create(newBookmark).Error
	if err != nil {
		return nil, dbutils.ParseDBError(err)
	}	
	return newBookmark, nil
}