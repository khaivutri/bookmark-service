package infrastructure

import (
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

const migrationPath = "file://./migrations"
func migrate(dbClient *gorm.DB) {
	

	err := sqldb.MigrateSQLDB(dbClient, migrationPath, "up", 0)
	if err != nil {
		panic(err)
	}
}