package bookmark

import (
	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/service/bookmark"
)
// BookmarkHandler defines the interface for handling bookmark-related HTTP requests.
type BookmarkHandler interface {
	// AddBookmark handles the creation of a new bookmark.
	AddBookmark(ctx *gin.Context)
	// GetBookmarks retrieves a paginated list of bookmarks.
	GetBookmarks(ctx *gin.Context)
	// UpdateBookmark updates an existing bookmark.
	UpdateBookmark(ctx *gin.Context) 
	// DeleteBookmark deletes a bookmark by its ID.
	DeleteBookmark(ctx *gin.Context)
}

type bookmarkHandler struct {
	svc bookmark.BookmarkService 
}

// NewBookmarkHandler constructs a new BookmarkHandler instance.
func NewBookmarkHandler(svc bookmark.BookmarkService) BookmarkHandler {
	return &bookmarkHandler{
		svc: svc,
	}
}