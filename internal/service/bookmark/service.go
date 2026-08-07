package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository/bookmark"

	"github.com/khaivutri/bookmark-service/pkg/utils"
)

// BookmarkService abstracts the business logic for managing bookmarks.
type BookmarkService interface {
	// AddBookmark generates a unique code, constructs a new bookmark, and persists it.
	AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error)

	// GetBookmarks retrieves a paginated list of bookmarks for a specific user.
	GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarkResponse, error)

	// UpdateBookmark updates description and URL values of an existing bookmark.
	UpdateBookmark(ctx context.Context, userID, bookmarkID, description, url string) error
	
	// DeleteBookmark deletes a bookmark by its unique ID.
	DeleteBookmark(ctx context.Context, userID, bookmarkID string) error
}

type bookmarkService struct {
	repo    bookmark.BookmarkRepository
	codeGen utils.GenCode
}

// NewBookmarkService constructs a new BookmarkService instance.
func NewBookmarkService(repo bookmark.BookmarkRepository, codeGen utils.GenCode) BookmarkService {
	return &bookmarkService{
		repo:    repo,
		codeGen: codeGen,
	}
}
