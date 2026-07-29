package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/docs"
	"github.com/khaivutri/bookmark-service/internal/api/middleware"
	"github.com/khaivutri/bookmark-service/internal/handler"
	v1 "github.com/khaivutri/bookmark-service/internal/handler/v1"
	userHandler "github.com/khaivutri/bookmark-service/internal/handler/v1/user"
	"github.com/khaivutri/bookmark-service/internal/repository"
	userRepo "github.com/khaivutri/bookmark-service/internal/repository/user"
	"github.com/khaivutri/bookmark-service/internal/service"
	userSvc "github.com/khaivutri/bookmark-service/internal/service/user"
	"github.com/khaivutri/bookmark-service/pkg/jwtutils"
	"github.com/khaivutri/bookmark-service/pkg/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	_ "github.com/khaivutri/bookmark-service/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Engine defines the interface for running and serving the API engine.
type Engine interface {
	// Start runs the HTTP server listening on the configured application port.
	Start() error
	// ServeHTTP routes HTTP requests to the underlying router (useful for testing).
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app 		*gin.Engine
	cfg 		*Config

	redis 		*redis.Client
	db 			*gorm.DB

	jwtGen    	jwtutils.JWTGenerator
	jwtVal   	jwtutils.JWTValidator
}

type EngineOpts struct {
	App 		*gin.Engine
	Cfg 		*Config

	Redis 		*redis.Client
	DB 			*gorm.DB

	JWTGen    	jwtutils.JWTGenerator
	JWTVal   	jwtutils.JWTValidator
}
// NewEngine initializes and returns a new Engine instance with defined routes and handlers.
func NewEngine(opts EngineOpts) Engine{
	app := &engine{
		app : 		opts.App,
		cfg : 		opts.Cfg,
		redis : 	opts.Redis,
		db: 		opts.DB,
		jwtGen: 	opts.JWTGen,
		jwtVal: 	opts.JWTVal,
	
	}
	app.initRoutes()
	return app
}

// Start runs the HTTP server listening on the configured application port.
func (e *engine) Start() error {
	return e.app.Run(fmt.Sprintf(":%s", e.cfg.AppPort))
}

// ServeHTTP routes HTTP requests to the underlying Gin engine.
func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

type handlers struct {
	healthCheckHandlers 	handler.HealthCheck
	linkHandler 			v1.ShortenURL	
	registerHandler 		userHandler.Handler
}

func (e *engine) initHandlers() *handlers {
	redisPinger := repository.NewRedisPinger(e.redis)
	dbPinger := repository.NewDBPinger(e.db)
	healthCheckSvc := service.NewHealthCheck(e.cfg.ServiceName, e.cfg.InstanceId, redisPinger, dbPinger)
	healthCheck := handler.NewHealthCheck(healthCheckSvc)

	urlStorage := repository.NewURLStorage(e.redis)
	shortenURLSvc := service.NewURLStorage(urlStorage, utils.NewGenCode())
	shortenURL := v1.NewShortenURL(shortenURLSvc)

	userRepo := userRepo.NewSqlRepository(e.db)
	hasher := utils.NewHasher()
	userSvc := userSvc.NewService(userRepo, hasher, e.jwtGen)
	registerHandler := userHandler.NewHandler(userSvc)
	return &handlers{healthCheck, shortenURL, registerHandler}
}

func (e *engine) initRoutes(){
	allHandlers := e.initHandlers()
	e.app.HandleMethodNotAllowed = true	
	
	// init middlewares
	jwtAuth := middleware.NewJWTAuth(e.jwtVal)
	
	docs.SwaggerInfo.BasePath = e.cfg.BasePath
	e.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	e.app.GET("/health-check", allHandlers.healthCheckHandlers.HealthCheck)

	v1 := e.app.Group("/v1") 
	{
		links := v1.Group("/links")
		{
			links.POST("/shorten", allHandlers.linkHandler.CreateShortenLink)
			links.GET("/redirect/:code", allHandlers.linkHandler.Redirect)
		}

		users := v1.Group("/users")
		{
			users.POST("/register", allHandlers.registerHandler.Register)
			users.POST("/login", allHandlers.registerHandler.Login)
		}

		self := v1.Group("/self")
		{	
			self.Use(jwtAuth.JWTAuth())
			self.GET("/info", allHandlers.registerHandler.GetSelfInfo)
			self.PUT("/info", allHandlers.registerHandler.UpdateSelfInfo)
		}
	}
}

