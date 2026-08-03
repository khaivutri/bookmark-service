package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository/bookmark"

	"github.com/khaivutri/bookmark-service/pkg/utils"
)

type BookmarkService interface {
	AddBookmark(ctx context.Context, descriptio, url, userID string) (*model.Bookmark, error)

	GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarkResponse, error)

	UpdateBookmark(ctx context.Context, bookmarkID, description, url string) error
	
	DeleteBookmark(ctx context.Context, bookmarkID string) error
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
