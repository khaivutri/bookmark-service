package sqldb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewClient initializes and returns a GORM database client based on environment configurations.
func NewClient(envPrefix string) (*gorm.DB, error) {
	cfg, err := newConfig(envPrefix)
	if err != nil {
		return nil, err
	}

	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}
	
	return db, nil
}