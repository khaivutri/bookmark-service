package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	userDTO "github.com/khaivutri/bookmark-service/internal/handler/v1/dto/user"
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


// UpdateSelfInfo update your current information
// @Summary      update your current information
// @Description  update your current information
// @Tags         User
// @Security     BearerAuth
// @Accept       application/json
// @Produce      application/json
// @Param        request  body      user.UpdateUserInfoReq  true  "Update user info request"
// @Success      200 {object} object{user.UpdateUserInfoRes} "Success"
// @Router       /v1/self/info [put]
func (u *userHandler) UpdateSelfInfo(ctx *gin.Context) {
	requestBody, err := requestutils.BindInputFromResquest[userDTO.UpdateUserInfoReq](ctx)
	if err != nil {
		return 
	}

	uid, err := requestutils.GetUserIDFromRequest(ctx)
	if err != nil {
		return 
	}

	if err = u.svc.UpdateSelfInfo(ctx, uid, requestBody.DisplayName, requestBody.Email); err != nil {
		log.Error().Err(err).Str("from", "handler.user.UpdateSelfInfo").Msg("failed to update self info")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
		return
	}
	
	response := userDTO.UpdateInfoRes{
		Message: "Edit current user successfully!",
	}
	ctx.JSON(http.StatusOK, response)
	
}