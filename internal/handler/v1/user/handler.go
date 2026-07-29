package user

import (
	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/service/user"
)

// Handler defines the interface for handling user-related HTTP requests.
type Handler interface {
	// Register processes user registration requests.
	Register(ctx *gin.Context)
	// Login processes user login requests.
	Login(ctx *gin.Context)
	// GetSelfInfo retrieves the profile of the currently authenticated user.
	GetSelfInfo(ctx *gin.Context)
	// UpdateSelfInfo updates the profile of the currently authenticated user.
	UpdateSelfInfo(ctx *gin.Context)
}

type userHandler struct {
	svc user.Service
}

// NewHandler constructs a new user handler instance.
func NewHandler(svc user.Service) Handler {
	return &userHandler{svc: svc}
}
