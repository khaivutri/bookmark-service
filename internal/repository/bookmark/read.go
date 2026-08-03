package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
)

func (b *bookmarkRepo) GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	bookmark := &model.Bookmark{}
	err := b.db.WithContext(ctx).Where("code = ?", code).First(bookmark).Error
	if err != nil {
		return nil, dbutils.ErrRecordNotFound
	}
	return bookmark, nil
}