package bookmark

import (
	"context"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
)

func (s *bookmarkService) AddBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	//create a new bookmark model
	code, errGen := s.codeGen.Generate(10)
	if errGen != nil {
		return nil, errGen
	}

	// check if code is unique
	record, errGet := s.repo.GetBookmarkByCode(ctx, code)
	if errGet != nil && errGet != dbutils.ErrRecordNotFound {
		return nil, errGet
	}

	// if code is not unique, generate a new code and try again
	if record != nil {
		return s.AddBookmark(ctx, description, url, userID)
	}

	bookmark := &model.Bookmark{
		Description: description,
		URL:         url,
		UserID:      userID,
		Code:        code,
	}


	// call repo 
	response, err := s.repo.CreateBookmark(ctx, bookmark)
	if err != nil {
		return nil, err
	}


	// return 
	return response, nil

}