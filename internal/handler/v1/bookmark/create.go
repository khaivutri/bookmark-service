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

// AddBookmark handles the creation of a new bookmark for a user.
// @Summary      Add Bookmark
// @Description  Creates a new bookmark for the authenticated user with the provided description and URL.
// @Tags         Bookmarks
// @Security     BearerAuth
// @Accept       application/json
// @Produce      application/json
// @Param        request  body      bookmark.BookmarkRequest  true  "Bookmark request"
// @Success      200      {object}  model.Bookmark            "Bookmark created successfully"
// @Failure      400      {object}  response.ErrMessage       "Bad Request / Invalid input"
// @Failure      401      {object}  response.ErrMessage       "Unauthorized"
// @Failure      409      {object}  response.ErrMessage       "Conflict / Duplicate bookmark code"
// @Failure      500      {object}  response.ErrMessage       "Internal Server Error"
// @Router       /v1/bookmarks [post]
func (h *bookmarkHandler) AddBookmark(ctx *gin.Context) {
	// Bind the incoming JSON request to the BookmarkRequest struct
	body, err := requestutils.BindInputFromResquest[bookmark.BookmarkRequest](ctx)
	if err != nil {
		return
	}

	// Get the user ID from the Gin context
	uid, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return
	}

	// call service
	res, err := h.svc.AddBookmark(ctx, body.Description, body.URL, uid)
	switch {
	case errors.Is(err, dbutils.ErrDuplicateBookmarkCode):
		ctx.AbortWithStatusJSON(http.StatusConflict,&response.ErrMessage{
			Message: "bookmark code invalid",
		})
		return
	case errors.Is(err, nil):
		break
	default:
		log.Error().Err(err).Str("from", "handler.bookmark.AddBookmark").Msg("failed to add bookmark")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
		return
	}

	ctx.JSON(200, res)
}