package main

import (
	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/api"
	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/jwtutils"
	"github.com/khaivutri/bookmark-service/pkg/logger"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/khaivutri/bookmark-service/pkg/sqldb"
	"github.com/khaivutri/bookmark-service/pkg/validation"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

//	@title Bookmark Service API
//	@version 2.0
//	@description This is a simple REST API for a bookmark service - Demo MLIoT Lab.
//	@BasePath /
//	@securityDefinitions.apikey BearerAuth
//	@in header
//	@name Authorization

func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	//set up validtation
	if err := validation.RegisterValidation(); err != nil {
		panic(err)
	}
	//set log level 
	logger.SetLogLevel(cfg.LogLevel)

	//set redis client
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}
	//set database
	dbClient, err := sqldb.NewClient("")
	if err != nil {
		panic(err)
	}
	err = dbClient.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	
	engine := createAPIApp(cfg, redisClient, dbClient)
	
	err = engine.Start()
	if err != nil {
		panic(err)
	}
}

func createAPIApp(cfg *api.Config, redis *redis.Client, db *gorm.DB) api.Engine{
	app := gin.New()

	jwtGen, err := jwtutils.NewJWTGenerator("./private.pem")
	if err != nil {
		panic(err)
	}
	jwtVal, err := jwtutils.NewJWTValidator("./public.pem")
	if err != nil {
		panic(err)
	}
	a := api.NewEngine(api.EngineOpts{
		App: 		app,
		Cfg: 		cfg,
		Redis: 		redis,
		DB: 		db,
		JWTGen: 	jwtGen,
		JWTVal: 	jwtVal,
	})
	return a
}	