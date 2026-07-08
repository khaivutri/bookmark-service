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
//@Failure 503 {object} model.HealthReport
//@Router /health-check [get]
func (hc *healthCheck) HealthCheck(ctx *gin.Context) {
	report, err := hc.IsHealthy.Check(ctx.Request.Context())

	if report == nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// err mirrors what's already captured in report.Message / report.Dependencies
	// (e.g. "DEGRADED" + dependency: {"redis": "DOWN"}). We still return the
	// full report to the caller instead of masking it behind a generic 500 -
	// hook up logging here if/when a logger is wired into this handler.
	if err != nil {
		_ = err
	}

	status := http.StatusOK
	if report.Message != "OK" {
		status = http.StatusServiceUnavailable
	}

	ctx.JSON(status, report)
}