package bookmark

import (
	"context"
)

// UpdateBookmark updates description and URL values of an existing bookmark.
func (s *bookmarkService) UpdateBookmark(ctx context.Context, bookmarkID, description, url string) error {
	err := s.repo.UpdateBookmark(ctx, bookmarkID, description, url)
	if err != nil {
		return err
	}
	return nil
}
