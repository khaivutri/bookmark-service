package user

import "gorm.io/gorm"

type sqlRepository struct {
	db *gorm.DB
}

// NewSqlRepository constructs a user Repository instance using a GORM database client.
func NewSqlRepository(db *gorm.DB) Repository {
	return &sqlRepository{db: db}
}
