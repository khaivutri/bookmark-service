package fixture

import (
	"github.com/khaivutri/bookmark-service/internal/model"
	"gorm.io/gorm"
)

type UserCommonTest struct {
	base // composition
}

func (u *UserCommonTest) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

func (u *UserCommonTest) GenerateData() error {
	// skip BeforeCreate
	db := u.db.Session(&gorm.Session{
		SkipHooks: true,
	})

	users := []*model.User{
		{
			Base:        GetTestBase("b649b57b-b7b6-44e4-a233-74147ecf56ee"),
			DisplayName: "Test1",
			UserName:    "test1",
			Password:    "pwd123",
			Email:       "test1@example.com",
		},
		{
			Base:        GetTestBase("b649b57b-b7b6-44e4-a233-74147ecf56ef"),
			DisplayName: "Test2",
			UserName:    "test2",
			Password:    "pwd123",
			Email:       "test2@example.com",
		},
	}

	return db.CreateInBatches(users, 10).Error
}
