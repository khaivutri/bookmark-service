package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/service"
)

// HealthCheck defines the interface for health check handler.
type HealthCheck interface {
	HealthCheck(ctx *gin.Context)
}

type healthCheck struct {
	IsHealthy service.HealthCheck
}

// NewHealthCheck creates and returns a new HealthCheck handler.
func NewHealthCheck(isHealthy service.HealthCheck) *healthCheck {
	return &healthCheck{IsHealthy: isHealthy}
}

// HealthCheck handles the health check HTTP request.
func (hc *healthCheck) HealthCheck(ctx *gin.Context) {
	report := hc.IsHealthy.Check()

	ctx.JSON(http.StatusOK, report)
}


