package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type Repository interface {
	CreateBookmark(ctx context.Context, newBookmark *model.Bookmark) (*model.Bookmark, error)
}

type bookmarkRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &bookmarkRepo{db: db}
}