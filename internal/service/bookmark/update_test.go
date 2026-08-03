package bookmark

import (
	"errors"
	"reflect"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	repoMocks "github.com/khaivutri/bookmark-service/internal/repository/bookmark/mocks"
	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"
	"github.com/stretchr/testify/mock"
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
	existing *model.Bookmark,
	getErr error,
	expected *model.Bookmark,
	updateErr error,
) *repoMocks.BookmarkRepository {
	t.Helper()

	repo := repoMocks.NewBookmarkRepository(t)
	repo.On("GetBookmarkByID", t.Context(), bookmarkID).
		Return(existing, getErr).
		Once()

	if getErr == nil {
		repo.On("UpdateBookmark", t.Context(), mock.MatchedBy(func(got *model.Bookmark) bool {
			return reflect.DeepEqual(got, expected)
		})).
			Return(updateErr).
			Once()
	}

	return repo
}

func testBookmark(id, description, url, code, userID string) *model.Bookmark {
	return &model.Bookmark{
		Base:        fixture.GetTestBase(id),
		Description: description,
		URL:         url,
		Code:        code,
		UserID:      userID,
	}
}

func TestBookmarkService_UpdateBookmark(t *testing.T) {
	t.Parallel()

	testcases := []updateBookmarkTestCase{
		{
			name: "success - updates description and URL, keeps other fields",
			setupMockRepo: func(t *testing.T) *repoMocks.BookmarkRepository {
				existing := testBookmark("bookmark-1", "Old description", "https://old.example.com", "code-1", "user-1")
				expected := testBookmark("bookmark-1", "New description", "https://new.example.com", "code-1", "user-1")
				return setupUpdateBookmarkRepo(t, existing.ID, existing, nil, expected, nil)
			},
			inputBookmarkID:  "bookmark-1",
			inputDescription: "New description",
			inputURL:         "https://new.example.com",
		},
		{
			name: "get bookmark fails - returns ErrFailGetBookmark",
			setupMockRepo: func(t *testing.T) *repoMocks.BookmarkRepository {
				return setupUpdateBookmarkRepo(t, "unknown-id", nil, errors.New("record not found"), nil, nil)
			},
			inputBookmarkID:  "unknown-id",
			inputDescription: "New description",
			inputURL:         "https://new.example.com",
			expectedErr:      ErrFailGetBookmark,
		},
		{
			name: "repository update fails - returns ErrFailUpdateBookmark",
			setupMockRepo: func(t *testing.T) *repoMocks.BookmarkRepository {
				existing := testBookmark("bookmark-3", "Old description", "https://old.example.com", "code-3", "user-3")
				expected := testBookmark("bookmark-3", "New description", "https://new.example.com", "code-3", "user-3")
				return setupUpdateBookmarkRepo(t, existing.ID, existing, nil, expected, errors.New("db write error"))
			},
			inputBookmarkID:  "bookmark-3",
			inputDescription: "New description",
			inputURL:         "https://new.example.com",
			expectedErr:      ErrFailUpdateBookmark,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) { runUpdateBookmarkTest(t, tc) })
	}
}
