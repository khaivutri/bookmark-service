package bookmark

import (
	"context"
	"errors"

	"github.com/khaivutri/bookmark-service/internal/model"
)

var (
	ErrCount 			= errors.New("failed to count bookmarks")
	ErrGetBookMarks 	= errors.New("failed to get bookmarks")
)

// GetBookmarks retrieves a list of bookmarks with limit and offset, along with the total count.
func (b *bookmarkRepo) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error) {
	bookmarks, err := b.getBookmarks(ctx, userID, "created_at ASC", limit, offset)
	if err != nil {
		return nil, 0, ErrGetBookMarks
	}

	count, err := b.countBookmarks(ctx, userID)
	if err != nil {
		return nil, 0, ErrCount
	}

	return bookmarks, count, nil
}



func (b *bookmarkRepo) getBookmarks(ctx context.Context, userID, sort string, limit, offset int) ([]*model.Bookmark, error) {
	bookmarks := make([]*model.Bookmark,limit)

	err := b.db.WithContext(ctx).Where("user_id = ?", userID).Order(sort).Limit(limit).Offset(offset).Find(&bookmarks).Error
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}


func (b *bookmarkRepo) countBookmarks(ctx context.Context, userID string) (int64, error) {
	var count int64

	err := b.db.WithContext(ctx).Model(&model.Bookmark{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}