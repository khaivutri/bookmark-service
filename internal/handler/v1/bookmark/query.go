package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/handler/v1/dto/bookmark"
	"github.com/khaivutri/bookmark-service/pkg/requestutils"
	"github.com/khaivutri/bookmark-service/pkg/response"
	"github.com/rs/zerolog/log"
)

type getBookmarkRequest struct {
	Page  int `form:"page" validate:"gte=1"`
	Limit int `form:"limit" validate:"gte=1"`
}

// GetBookmarks handles retrieving a paginated list of bookmarks for the authenticated user.
// @Summary      Get user bookmarks
// @Description  Retrieves all bookmarks belonging to the authenticated user with pagination.
// @Tags         Bookmarks
// @Security     BearerAuth
// @Produce      application/json
// @Param        page   query     int                                             false  "Page number"    default(1)  minimum(1)
// @Param        limit  query     int                                             false  "Items per page" default(10) minimum(1)
// @Success      200    {object}  bookmark.GetBookmarksResponse                  "Bookmarks retrieved successfully"
// @Failure      400    {object}  response.ErrMessage                             "Bad Request / Invalid input"
// @Failure      401    {object}  response.ErrMessage                             "Unauthorized"
// @Failure      500    {object}  response.ErrMessage                             "Internal Server Error"
// @Router       /v1/bookmarks [get]
func (b *bookmarkHandler) GetBookmarks(ctx *gin.Context) {
	queryParam, err := requestutils.BindInputFromResquest[getBookmarkRequest](ctx)
	if err != nil {
		return
	}

	uid, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return
	}

	res, err := b.svc.GetBookmarks(ctx, uid, queryParam.Page, queryParam.Limit)

	switch {
	case errors.Is(err, nil):
		break
	default:
		log.Error().Err(err).Str("from", "handler.bookmark.GetBookmarks").Msg("failed to get bookmarks")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, &bookmark.GetBookmarksResponse{
		Data: res.Bookmarks,
		Pagination: &bookmark.Pagination{
			Page:  queryParam.Page,
			Limit: queryParam.Limit,
			Total: res.Total,
		},
	})

}
