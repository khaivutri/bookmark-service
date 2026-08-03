package bookmark

import (
	"context"

)

// DeleteBookmark deletes a bookmark by its unique ID.
func (svc *bookmarkService) DeleteBookmark(ctx context.Context, bookmarkID string) error {
	err := svc.repo.DeleteBookmarkByID(ctx, bookmarkID)
	if err != nil {
		return err
	}
	return nil
}