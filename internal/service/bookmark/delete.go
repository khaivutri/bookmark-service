package bookmark

import (
	"context"
)

// DeleteBookmark deletes a bookmark by its unique ID.
func (svc *bookmarkService) DeleteBookmark(ctx context.Context, userID, bookmarkID string) error {
	err := svc.repo.DeleteBookmarkByID(ctx, userID, bookmarkID)
	if err != nil {
		return err
	}
	return nil
}
