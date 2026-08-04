package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
)

// UpdateBookmark saves changes to an existing bookmark record.
func (r *bookmarkRepo) UpdateBookmark(ctx context.Context, bookmarkID, description, url string) error {
	res := r.db.WithContext(ctx).Where("id = ?", bookmarkID).Updates(model.Bookmark{
		Description: description,
		URL:         url,
	})

	if res.Error != nil {
		return dbutils.ParseDBError(res.Error)
	}

	// GORM's Updates returns RowsAffected == 0 if no row matches the WHERE condition
	if res.RowsAffected == 0 {
		return dbutils.ErrRecordNotFound
	}

	return nil
}
