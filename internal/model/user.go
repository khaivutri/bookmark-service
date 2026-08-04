package model


type User struct {
	Base
	DisplayName string 			`json:"display_name" gorm:"not null"`
	UserName    string 			`json:"username" gorm:"column:username;uniqueIndex:uni_users_username"`
	Password    string 			`json:"-"`
	Email       string 			`json:"email" gorm:"unique"`
}

