package user

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20" example:"johndoe"`
	Password string `json:"password" binding:"required,password" example:"Password123@"`
}

