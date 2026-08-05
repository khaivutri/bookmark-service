package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	userDTO "github.com/khaivutri/bookmark-service/internal/handler/v1/dto/user"
	"github.com/khaivutri/bookmark-service/pkg/dbutils"
	"github.com/khaivutri/bookmark-service/pkg/requestutils"
	"github.com/khaivutri/bookmark-service/pkg/response"
	"github.com/rs/zerolog/log"
)

// @Summary      Register a new user
// @Description  Creates a new user account with the provided username, display name, password, and email.
// @Tags         User
// @Accept       application/json
// @Produce      application/json
// @Param 		 request body user.RegisterRequest true "Register request"
// @Success      201  {object}  user.RegisterResponse  "Register response"
// @Failure      400  {object}  map[string]string  "Bad Request"
// @Failure      409  {object}  map[string]string  "Conflict"
// @Failure      500  {object}  map[string]string  "Internal Server Error"
// @Router       /v1/users/register [post]
// Register handles user registration requests.
func (u *userHandler) Register(ctx *gin.Context) {
	body, err := requestutils.BindInputFromResquest[userDTO.RegisterRequest](ctx)
	if err != nil {
		return
	}

	user, err := u.svc.CreateUser(ctx, body.Username, body.DisplayName, body.Password, body.Email)
	switch {
	case errors.Is(err, dbutils.ErrDuplicateUserName):
		ctx.AbortWithStatusJSON(http.StatusConflict, response.ErrMessage{
			Message: "username already exists",
		})
		return
	case errors.Is(err, dbutils.ErrDuplicateEmail):
		ctx.AbortWithStatusJSON(http.StatusConflict, response.ErrMessage{
			Message: "email already exists",
		})
		return
	case err == nil:
	default:
		log.Error().Err(err).Str("from", "handler.user.Register").Msg("failed to create user")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, response.InternalServerErrResponse)
		return
	}
	response := &userDTO.RegisterResponse{
		Data: userDTO.RegisterData{
			ID: user.ID,
			Username: user.UserName,
			DisplayName: user.DisplayName,
			Email: user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Message: "User registered successfully!",
	}
	ctx.JSON(http.StatusCreated, response)
}
