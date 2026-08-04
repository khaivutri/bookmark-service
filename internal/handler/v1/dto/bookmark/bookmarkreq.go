package bookmark

type BookmarkRequest struct {
	Description 	string 	`json:"description" exmple"Google" validate:"lte=255"`
	URL		 		string 	`json:"url" example:"https://www.google.com" validate:"required,url,lte=2048"`
}