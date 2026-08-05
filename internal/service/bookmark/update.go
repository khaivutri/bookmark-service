package bookmark

import (
	"context"
)

// UpdateBookmark updates description and URL values of an existing bookmark.
func (s *bookmarkService) UpdateBookmark(ctx context.Context, userID, bookmarkID, description, url string) error {
	err := s.repo.UpdateBookmark(ctx, userID, bookmarkID, description, url)
	if err != nil {
		return err
	}
	return nil
}
