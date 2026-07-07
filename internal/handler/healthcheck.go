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
func NewHealthCheck(isHealthy service.HealthCheck) HealthCheck {
	return &healthCheck{IsHealthy: isHealthy}
}

// HealthCheck handles the health check HTTP request.
//@Summary Health Check
//@Tags Health Check
//@Accept json
//@Produce json
//@Success 200 {object} model.HealthReport	
//@Router /health-check [get]
func (hc *healthCheck) HealthCheck(ctx *gin.Context) {
	report := hc.IsHealthy.Check()

	ctx.JSON(http.StatusOK, report)
}


