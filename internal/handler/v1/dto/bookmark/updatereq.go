package bookmark



type UpdateRequest struct {
	ID          string `form:"id" validate:"required"`
	Description string `json:"description" example:"Google" validate:"lte=255"`
	URL         string `json:"url" example:"https://www.google.com" validate:"required,url,lte=2048"`
}