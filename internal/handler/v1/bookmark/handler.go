package bookmark

import (
	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/service/bookmark"
)
type BookmarkHandler interface {
	AddBookmark(ctx *gin.Context)
	GetBookmarks(ctx *gin.Context) 
}

type bookmarkHandler struct {
	svc bookmark.BookmarkService 
}

func NewBookmarkHandler(svc bookmark.BookmarkService) BookmarkHandler {
	return &bookmarkHandler{
		svc: svc,
	}
}