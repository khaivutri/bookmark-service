package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
)

// DeleteBookmarkByID deletes a bookmark record by its unique ID.
func (b *bookmarkRepo) DeleteBookmarkByID(ctx context.Context, userID, bookmarkID string) error {
	res := b.db.WithContext(ctx).Where("id = ?", bookmarkID).Delete(&model.Bookmark{})
	if res.Error != nil {
		return dbutils.ParseDBError(res.Error)
	}

	// GORM's Delete does not return gorm.ErrRecordNotFound when no row matches
	// the WHERE condition — it only reports RowsAffected == 0 with a nil error.
	// So we must check RowsAffected ourselves and return the mapped error directly.
	if res.RowsAffected == 0 {
		return dbutils.ErrRecordNotFound 
	}

	return nil
}