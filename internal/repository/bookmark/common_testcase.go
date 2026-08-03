package bookmark

import (
	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/test/data/fixture"
)

var (
	existingBookmark1 = &model.Bookmark{
		Base:        fixture.GetTestBase("b649b57b-b7b6-44e4-a233-74147ecf56ee"),
		Description: "description1",
		URL:         "https://example.com",
		Code:        "code1",
		UserID:		 "b649b57b-b7b6-44e4-a233-74147ecf56ee",
	}
	existingBookmark2 = &model.Bookmark{
		Base:        fixture.GetTestBase("b649b57b-b7b6-44e4-a233-74147ecf56ef"),
		Description: "description2",
		URL:         "https://example.com",
		Code:        "code2",
		UserID:      "b649b57b-b7b6-44e4-a233-74147ecf56ee",
	}
)