package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository/bookmark"

	"github.com/khaivutri/bookmark-service/pkg/utils"
)

type BookmarkService interface {
	AddBookmark(ctx context.Context, descriptio, url, userID string) (*model.Bookmark, error)

	GetBookmarks(ctx context.Context, userID string, page, limit int) (*getBookmarkResponse, error)
}

type bookmarkService struct {
	repo    bookmark.BookmarkRepository
	codeGen utils.GenCode
}

func NewBookmarkService(repo bookmark.BookmarkRepository, codeGen utils.GenCode) BookmarkService {
	return &bookmarkService{
		repo:    repo,
		codeGen: codeGen,
	}
}

// GetBookmarkResponse is exported so generated service mocks can implement
// BookmarkService from another package.
type GetBookmarkResponse = getBookmarkResponse
