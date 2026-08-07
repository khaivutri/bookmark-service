package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"gorm.io/gorm"
)

// BookmarkRepository defines the data storage operations for bookmarks.
type BookmarkRepository interface {
	// CreateBookmark inserts a new bookmark record into the database.
	CreateBookmark(ctx context.Context, newBookmark *model.Bookmark) (*model.Bookmark, error)

	// GetBookmarks retrieves a list of bookmarks with limit and offset, along with the total count.
	GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error)

	// GetBookmarkByID queries a bookmark by its unique ID.
	GetBookmarkByID(ctx context.Context, id string) (*model.Bookmark, error)

	// GetBookmarkByCode queries a bookmark by its unique short code.
	GetBookmarkByCode(ctx context.Context, code string) (*model.Bookmark, error)

	// UpdateBookmark saves changes to an existing bookmark record.
	UpdateBookmark(ctx context.Context, userID, bookmarkID, description, url string) error

	// DeleteBookmarkByID deletes a bookmark record by its unique ID.
	DeleteBookmarkByID(ctx context.Context, userID, bookmarkID string) error
}

type bookmarkRepo struct {
	db *gorm.DB
}

// NewRepository constructs a new BookmarkRepository instance using a GORM database client.
func NewRepository(db *gorm.DB) BookmarkRepository {
	return &bookmarkRepo{db: db}
}
