package sqldb

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func MigrateSQLDB(db *gorm.DB, migrationPath, mode string, steps int) error {
	sqldb, err := db.DB()
	if err != nil {
		return err
	}

	pgDriver, err := postgres.WithInstance(sqldb, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationPath, db.Name(), pgDriver)
	if err != nil {  
		return err
	}

	return migrateSchema(m, mode, steps)
}

func migrateSchema(m *migrate.Migrate, mode string, steps int) error {
	var migrationErr error
	
	switch mode {
		case "up":
			migrationErr = m.Up()
		case "down":
			migrationErr = m.Down()
		case "steps":
			if steps ==0 {
				return errors.New("[DB Migration] Steps must not be 0. Please provide a positive number.")
			}
			migrationErr = m.Steps(steps)
		default:
			return errors.New("[DB Migration] Invalid mode. Please use either 'up', 'down' or 'step' as mode.")
	}

	if migrationErr != nil && migrationErr != migrate.ErrNoChange {
		return fmt.Errorf("[DB Migration] Error: %s", migrationErr.Error())
	}

	return nil
}