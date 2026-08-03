package bookmark

import (
	"context"

)

func (svc *bookmarkService) DeleteBookmark(ctx context.Context, bookmarkID string) error {
	err := svc.repo.DeleteBookmarkByID(ctx, bookmarkID)
	if err != nil {
		return err
	}
	return nil
}