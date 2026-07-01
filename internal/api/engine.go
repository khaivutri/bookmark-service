package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/handler"
	"github.com/khaivutri/bookmark-service/internal/service"
)

// Engine defines the interface for the API engine.
type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *Config
}

// NewEngine creates and returns a new Engine instance with initialized routes.
func NewEngine(cfg *Config) Engine{
	app := &engine{
		app : gin.Default(),
		cfg : cfg,
	}
	app.initRoutes()
	return app
}

// Start runs the API server on the configured port.
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP serves the HTTP request using the underlying Gin engine.
func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}
	
func (e *engine) initRoutes(){
	healthCheckSvc := service.NewHealthCheck(e.cfg)
	healthCheck := handler.NewHealthCheck(healthCheckSvc)
	
	e.app.HandleMethodNotAllowed = true	
	e.app.GET("/health-check", healthCheck.HealthCheck)
}

