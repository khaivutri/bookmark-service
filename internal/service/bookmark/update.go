package bookmark

import (
	"context"
	"errors"

)

var (
	ErrFailGetBookmark    = errors.New("failed to get bookmark")
	ErrFailUpdateBookmark = errors.New("failed to update bookmark")

)
func (s *bookmarkService) UpdateBookmark(ctx context.Context, bookmarkID, description, url string) error {
	// get bookmark by id
	bookmark, err := s.repo.GetBookmarkByID(ctx, bookmarkID)
	if err != nil {
		return ErrFailGetBookmark 
	}

	// update bookmark
	bookmark.Description = description
	bookmark.URL = url

	err = s.repo.UpdateBookmark(ctx, bookmark)
	if err != nil {
		return ErrFailUpdateBookmark
	}

	return nil
}