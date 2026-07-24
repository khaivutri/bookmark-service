package repository

import (
	"context"

	"gorm.io/gorm"
)

// DBPinger is a wrapper around gorm.DB for health checking.
type DBPinger struct {
	db *gorm.DB
}

// NewDBPinger returns a new DBPinger.
func NewDBPinger(db *gorm.DB) *DBPinger {
	return &DBPinger{db: db}
}

// Ping verifies the database connection by retrieving the sql.DB and pinging it.
func (p *DBPinger) Ping(ctx context.Context) error {
	if p.db == nil {
		return ErrDependencyDown
	}

	sqlDB, err :=p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}
