package bookmark

type DeleteRequest struct {
	ID string `form:"id" validate:"required"`
}