package user

import "time"

type RegisterData struct {
	ID          	string    		`json:"id"`
	Username    	string    		`json:"username"`
	Email       	string    		`json:"email"`
	DisplayName 	string    		`json:"display_name"`
	CreatedAt   	time.Time 		`json:"created_at"`
	UpdatedAt   	time.Time 		`json:"updated_at"`
}
type RegisterResponse struct {
	Data RegisterData `json:"data"`
	Message string `json:"message"`
}