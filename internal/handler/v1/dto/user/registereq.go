package user

type RegisterRequest struct {
	Username 		string 		`json:"username" validate:"required,min=3,max=16"`
	DisplayName 	string 		`json:"display_name" validate:"required,min=3,max=16"`
	Password 		string 		`json:"password" validate:"required,min=6,max=16"`
	Email    		string 		`json:"email" validate:"required,email"`
}