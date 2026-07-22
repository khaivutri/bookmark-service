package fixture

import (
	"testing"

	"github.com/khaivutri/bookmark-service/pkg/sqldb"
	"gorm.io/gorm"
)
type Fixture interface {
	SetupDB(db *gorm.DB)
	Migrate() 				error
	GenerateData() 			error

	DB() 					*gorm.DB
}

type base struct {
	db *gorm.DB
}

func (b *base) SetupDB(db *gorm.DB) {
	b.db = db
}

func (b *base) DB() *gorm.DB {
	return b.db
}

func NewFixture(t *testing.T, fixture Fixture) *gorm.DB {
	// create test db
	fixture.SetupDB(sqldb.InitMockDB(t))

	// migrare database schema
	err := fixture.Migrate()
	if err != nil {
		t.Fatalf("fail to migrate fixture: %s", err)
	}
	// generate data
	err = fixture.GenerateData()
	if err != nil {
		t.Fatalf("fail to generate data: %s", err)
	}

	// return db
	return fixture.DB()
}
