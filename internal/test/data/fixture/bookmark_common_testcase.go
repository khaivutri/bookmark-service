package fixture

import (
	"github.com/khaivutri/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type BookmarkCommonTestDB struct {
	UserCommonTest
}

func (b *BookmarkCommonTestDB) Migrate() error {
	return b.DB().AutoMigrate(&model.Bookmark{}, &model.User{})
}

func (b *BookmarkCommonTestDB) GenerateData() error {
	err := b.UserCommonTest.GenerateData()
	if err != nil {
		return err
	}

	bookmarks := []*model.Bookmark{ 
		{
			Base: GetTestBase("b649b57b-b7b6-44e4-a233-74147ecf56ee"),
			Description: "description1",
			URL: "https://example.com",
			Code: "code1",
			UserID: "b649b57b-b7b6-44e4-a233-74147ecf56ee",
		},
		{
			Base: GetTestBase("b649b57b-b7b6-44e4-a233-74147ecf56ef"),
			Description: "description2",
			URL: "https://example.com",
			Code: "code2",
			UserID: "b649b57b-b7b6-44e4-a233-74147ecf56ee",
		},
	}

	return b.db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(bookmarks, 10).Error
}