package user

type RegisterRequest struct {
    Username    string `json:"username" binding:"required,min=3,max=20" example:"johndoe"`
    DisplayName string `json:"display_name" binding:"required,min=3,max=20" example:"John Doe"`
    Password    string `json:"password" binding:"required,password" example:"Password123@"`
    Email       string `json:"email" binding:"required,email" example:"john.doe@example.com"`
}