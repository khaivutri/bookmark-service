package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
	bookmarkDTO "github.com/khaivutri/bookmark-service/internal/handler/v1/dto/bookmark"
	"github.com/khaivutri/bookmark-service/pkg/requestutils"
	"github.com/khaivutri/bookmark-service/pkg/response"
	"github.com/rs/zerolog/log"
)

// UpdateBookmark handles updating an existing bookmark for the authenticated user.
// @Summary      Update bookmark
// @Description  Updates the description and/or URL of an existing bookmark belonging to the authenticated user.
// @Tags         Bookmarks
// @Security     BearerAuth
// @Accept       application/json
// @Produce      application/json
// @Param        id       query     string                   true   "Bookmark ID"
// @Param        request  body      bookmark.UpdateRequest   true   "Update bookmark request payload"
// @Success      200      {object}  bookmark.UpdateResponse  "Bookmark updated successfully"
// @Failure      400      {object}  response.ErrMessage      "Bad Request / Invalid input"
// @Failure      401      {object}  response.ErrMessage      "Unauthorized"
// @Failure      500      {object}  response.ErrMessage      "Internal Server Error"
// @Router       /v1/bookmarks [put]
func (b *bookmarkHandler) UpdateBookmark(ctx *gin.Context) {
	// binding input
	input, err :=requestutils.BindInputFromResquest[bookmarkDTO.UpdateRequest](ctx)
	if err != nil {
		return
	}
	
	// call service to update bookmark
	err = b.svc.UpdateBookmark(ctx, input.ID, input.Description, input.URL)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.bookmark.UpdateBookmark").Msg("Fail to update bookmark")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrMessage{Message: "Fail to update bookmark"})
		return
	}
	ctx.JSON(http.StatusOK, bookmarkDTO.UpdateResponse{Message: "Success"})
}