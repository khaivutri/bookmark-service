package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type BookmarkRepository interface {
	CreateBookmark(ctx context.Context, newBookmark *model.Bookmark) (*model.Bookmark, error)

	GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error)
	GetBookmarkByID(ctx context.Context, id string) (*model.Bookmark, error)
	GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error)
	
	UpdateBookmark(ctx context.Context, bookmark *model.Bookmark) (error)

	DeleteBookmarkByID(ctx context.Context, id string) error
}

type bookmarkRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) BookmarkRepository {
	return &bookmarkRepo{db: db}
}