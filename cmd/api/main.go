package main

import (
	"github.com/khaivutri/bookmark-service/internal/api"
	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/pkg/logger"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
	"github.com/khaivutri/bookmark-service/pkg/sqldb"
	"github.com/khaivutri/bookmark-service/pkg/validation"
)

//@title Bookmark Service API
//@version 1.5
//@description This is a simple REST API for a bookmark service.
//@BasePath /
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
	
	engine := api.NewEngine(cfg, redisClient, dbClient)
	err = engine.Start()

	if err != nil {
		panic(err)
	}
}