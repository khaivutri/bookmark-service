package user

import (
	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/service/user"
)

type Handler interface {
	Register(ctx *gin.Context)
}

type userHandler struct {
	svc user.Service
}

func NewHandler(svc user.Service) Handler {
	return &userHandler{svc: svc}
}
