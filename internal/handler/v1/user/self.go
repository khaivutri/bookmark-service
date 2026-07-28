package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/pkg/requestutils"
	"github.com/khaivutri/bookmark-service/pkg/response"
	"github.com/rs/zerolog/log"
)

// GetSelfInfo   get your current information
// @Summary      get your current information
// @Description  get your current information
// @Tags         User
// @Security     BearerAuth
// @Accept       application/json
// @Produce      application/json
// @Success      200 {object} object{data=model.User} "Success"
// @Router       /v1/self/info [get]
func (u *userHandler) GetSelfInfo(ctx *gin.Context) {
	// get claims in context 
	uid, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return 
	}

	currUser, err := u.svc.GetSelfInfo(ctx, uid)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.user.Login").Msg("failed to get self info")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
		return
	}

	ctx.JSON(http.StatusOK, currUser)
}