package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, newBookmark *model.Bookmark) (*model.Bookmark, error)

	GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error)
	GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error)
}

type bookmarkRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) BookmarkRepository {
	return &bookmarkRepo{db: db}
}