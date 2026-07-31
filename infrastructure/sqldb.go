package infrastructure

import (
	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/sqldb"
	"gorm.io/gorm"
)

func createDB() *gorm.DB{
	dbClient, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}
	return dbClient
}

func migrate(dbClient *gorm.DB) {
	err := dbClient.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
}