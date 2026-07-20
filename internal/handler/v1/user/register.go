package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	userDTO "github.com/khaivutri/bookmark-service/internal/handler/v1/dto/user"
	"github.com/rs/zerolog/log"
)


func (u *userHandler) Register(ctx *gin.Context) {
	body := &userDTO.RegisterRequest{}

	if err := ctx.ShouldBind(&body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := u.svc.CreateUser(ctx, body.Username, body.DisplayName, body.Password, body.Email)
	if err != nil {
		log.Error().Err(err).Str("from", "handler.user.Register").Msg("failed to register user")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}