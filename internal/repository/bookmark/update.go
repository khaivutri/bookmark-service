package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
)

func (r *bookmarkRepo) UpdateBookmark(ctx context.Context, bookmark *model.Bookmark) (error) {
	err := r.db.WithContext(ctx).Save(bookmark).Error
	if err != nil {
		return dbutils.ParseDBError(err)
	}
	return nil
}