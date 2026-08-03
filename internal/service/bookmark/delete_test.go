package bookmark

import (
	"errors"
	"testing"

	repoMocks "github.com/khaivutri/bookmark-service/internal/repository/bookmark/mocks"
	"github.com/stretchr/testify/require"
)

func TestBookmarkService_DeleteBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		bookmarkID string
		repoErr    error
	}{
		{
			name:       "deletes bookmark successfully",
			bookmarkID: "bookmark-1",
		},
		{
			name:       "returns repository error",
			bookmarkID: "bookmark-2",
			repoErr:    errors.New("delete bookmark failed"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo := repoMocks.NewBookmarkRepository(t)
			repo.On("DeleteBookmarkByID", t.Context(), tc.bookmarkID).
				Return(tc.repoErr).
				Once()

			err := (&bookmarkService{repo: repo}).DeleteBookmark(t.Context(), tc.bookmarkID)

			require.ErrorIs(t, err, tc.repoErr)
		})
	}
}
