package user

type LoginResponse struct {
	Message string `json:"message" example:"Logged in successfully!"`
	Data 	string `json:"data" example:"jwt_token"`
}