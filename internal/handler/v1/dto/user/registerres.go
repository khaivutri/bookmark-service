package user

import "time"

type RegisterData struct {
	ID          	string    		`json:"id" example:"7d80f755-7dce-4c95-b8bf-75bb8e240ef2"`
	Username    	string    		`json:"username" example:"johndoe" example:"johndoe"`
	Email       	string    		`json:"email" example:"john.doe@example.com" example:"john.doe@example.com"`
	DisplayName 	string    		`json:"display_name" example:"John Doe"`
	CreatedAt   	time.Time 		`json:"created_at" example:"2026-07-21T12:34:56Z"`
	UpdatedAt   	time.Time 		`json:"updated_at" example:"2026-07-21T12:34:56Z"`
}
type RegisterResponse struct {
	Data RegisterData `json:"data"`
	Message string `json:"message"`
}