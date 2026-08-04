package bookmark

import (
	"errors"
	"testing"

	repoMocks "github.com/khaivutri/bookmark-service/internal/repository/bookmark/mocks"

	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/stretchr/testify/require"
)

type updateBookmarkTestCase struct {
	name             string
	setupMockRepo    func(t *testing.T) *repoMocks.BookmarkRepository
	inputBookmarkID  string
	inputDescription string
	inputURL         string
	expectedErr      error
}

var errUpdateBookmark = errors.New("db write error")

func runUpdateBookmarkTest(t *testing.T, tc updateBookmarkTestCase) {
	t.Helper()
	t.Parallel()

	repo := tc.setupMockRepo(t)
	svc := &bookmarkService{repo: repo}

	err := svc.UpdateBookmark(t.Context(), tc.inputBookmarkID, tc.inputDescription, tc.inputURL)

	if tc.expectedErr != nil {
		require.ErrorIs(t, err, tc.expectedErr)
		return
	}

	require.NoError(t, err)
}

func setupUpdateBookmarkRepo(
	t *testing.T,
	bookmarkID string,
	description string,
	url string,
	updateErr error,
) *repoMocks.BookmarkRepository {
	t.Helper()

	repo := repoMocks.NewBookmarkRepository(t)
	repo.On("UpdateBookmark", t.Context(), bookmarkID, description, url).
		Return(updateErr).
		Once()

	return repo
}

func TestBookmarkService_UpdateBookmark(t *testing.T) {
	t.Parallel()

	testcases := []updateBookmarkTestCase{
		{
			name: "success - updates description and URL",
			setupMockRepo: func(t *testing.T) *repoMocks.BookmarkRepository {
				return setupUpdateBookmarkRepo(t, "bookmark-1", "New description", "https://new.example.com", nil)
			},
			inputBookmarkID:  "bookmark-1",
			inputDescription: "New description",
			inputURL:         "https://new.example.com",
		},
		{
			name: "repository update fails - returns repository error",
			setupMockRepo: func(t *testing.T) *repoMocks.BookmarkRepository {
				return setupUpdateBookmarkRepo(t, "unknown-id", "New description", "https://new.example.com", dbutils.ErrRecordNotFound)
			},
			inputBookmarkID:  "unknown-id",
			inputDescription: "New description",
			inputURL:         "https://new.example.com",
			expectedErr:      dbutils.ErrRecordNotFound,
		},
		{
			name: "repository update returns database error",
			setupMockRepo: func(t *testing.T) *repoMocks.BookmarkRepository {
				return setupUpdateBookmarkRepo(t, "bookmark-3", "New description", "https://new.example.com", errUpdateBookmark)
			},
			inputBookmarkID:  "bookmark-3",
			inputDescription: "New description",
			inputURL:         "https://new.example.com",
			expectedErr:      errUpdateBookmark,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) { runUpdateBookmarkTest(t, tc) })
	}
}
