package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	userDTO "github.com/khaivutri/bookmark-service/internal/handler/v1/dto/user"
	"github.com/khaivutri/bookmark-service/internal/service/user"
	"github.com/khaivutri/bookmark-service/pkg/requestutils"
	"github.com/khaivutri/bookmark-service/pkg/response"
	"github.com/rs/zerolog/log"
)

// @Summary      User Login
// @Description  Logs in a user with the provided username and password, returning a JWT token.
// @Tags         User
// @Accept       application/json
// @Produce      application/json
// @Param        request  body      user.LoginRequest  true  "Login request"
// @Success      200  {object}  user.LoginResponse  "Login response"
// @Failure      400  {object}  response.ErrMessage  "Bad Request / Invalid credentials"
// @Failure      500  {object}  response.ErrMessage  "Internal Server Error"
// @Router       /v1/users/login [post]
func (u *userHandler) Login(ctx *gin.Context) {
	body, err := requestutils.BindInputFromResquest[userDTO.LoginRequest](ctx)
	if err != nil {
		return
	}

	token, err := u.svc.Login(ctx, body.Username, body.Password)
	switch {
		case errors.Is(err, user.ErrInvalidCredential):
			ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrMessage{Message: "invalid credential"})
			return
		case err == nil:
		default: 
			log.Error().Err(err).Str("from", "handler.user.Login").Msg("failed to login")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
			return
	}

	ctx.JSON(http.StatusOK, userDTO.LoginResponse{
													Message: 	"Logged in successfully!", 
													Data: 		token,
												})
}