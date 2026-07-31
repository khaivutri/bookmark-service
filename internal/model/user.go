package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          string    		`json:"id" gorm:"type:uuid;primaryKey"`
	DisplayName string 			`json:"display_name" gorm:"not null"`
	UserName    string 			`json:"username" gorm:"column:username;uniqueIndex:uni_users_username"`
	Password    string 			`json:"-"`
	Email       string 			`json:"email" gorm:"unique"`
	CreatedAt   time.Time 		`json:"created_at" gorm:"autoCreateTime"`
	DeletedAt 	gorm.DeletedAt 	`json:"-" gorm:"index"`
	UpdatedAt   time.Time 		`json:"updated_at" gorm:"autoUpdateTime"`
}

// BeforeCreate generates a new UUID for the user prior to inserting the record.
func (u *User) BeforeCreate(tx *gorm.DB) ( err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return 
}