package main

import (
	"github.com/khaivutri/bookmark-service/internal/api"
	redisPkg "github.com/khaivutri/bookmark-service/pkg/redis"
)

//@title Bookmark Service API
//@version 1.0
//@description This is a simple REST API for a bookmark service.
//@BasePath /
//@host localhost:8080
func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}
	
	redisClient, err := redisPkg.NewClient("")
	if err != nil {
		panic(err)
	}
	engine := api.NewEngine(cfg, redisClient)
	err = engine.Start()

	if err != nil {
		panic(err)
	}
}