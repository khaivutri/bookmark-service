package bookmark

import (
	"context"
	"github.com/khaivutri/bookmark-service/internal/model"
)


type GetBookmarkResponse struct {
	Bookmarks 	[]*model.Bookmark 	`json:"bookmarks"`
	Total		int64				`json:"total"`	
}

func (s *bookmarkService) GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarkResponse, error) {
	offset := (page -1 )*limit

	bookmarks, count, err := s.repo.GetBookmarks(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &GetBookmarkResponse{	
									Bookmarks: bookmarks, 
									Total: count,
								}, nil
}
