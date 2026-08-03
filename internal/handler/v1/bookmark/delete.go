package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/handler/v1/dto/bookmark"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/khaivutri/bookmark-service/pkg/requestutils"
	"github.com/khaivutri/bookmark-service/pkg/response"
	"github.com/rs/zerolog/log"
)

// DeleteBookmark handles deleting an existing bookmark for the authenticated user.
// @Summary      Delete bookmark
// @Description  Deletes a bookmark belonging to the authenticated user by its ID.
// @Tags         Bookmarks
// @Security     BearerAuth
// @Produce      application/json
// @Param        id  query     string                   true  "Bookmark ID"
// @Success      200  {object}  bookmark.DeleteResponse  "Bookmark deleted successfully"
// @Failure      400  {object}  response.ErrMessage      "Bad Request / Invalid input"
// @Failure      401  {object}  response.ErrMessage      "Unauthorized"
// @Failure      404  {object}  response.ErrMessage      "Bookmark not found"
// @Failure      500  {object}  response.ErrMessage      "Internal Server Error"
// @Router       /v1/bookmarks [delete]
func (b *bookmarkHandler) DeleteBookmark(ctx *gin.Context) {
	query, err := requestutils.BindInputFromResquest[bookmark.DeleteRequest](ctx)
	if err != nil {
		return
	}

	err = b.svc.DeleteBookmark(ctx, query.ID)
	switch {
	case errors.Is(err, nil):
		break
	case errors.Is(err, dbutils.ErrRecordNotFound):
		ctx.AbortWithStatusJSON(http.StatusNotFound, response.ErrMessage{Message: "Bookmark not found"})
		return
	default:
		log.Error().Err(err).Str("from", "handler.bookmark.DeleteBookmark").Msg("Fail to delete bookmark")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrMessage{Message: "Fail to delete bookmark"})
		return
	}

	ctx.JSON(http.StatusOK, bookmark.DeleteResponse{Message: "Success"})

}
