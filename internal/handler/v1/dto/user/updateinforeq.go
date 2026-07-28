package user

type UpdateUserInfoReq struct {
	DisplayName string `json:"display_name" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
}