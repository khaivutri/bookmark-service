package main

import "github.com/khaivutri/bookmark-service/internal/api"

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
	
	engine := api.NewEngine(cfg)
	err = engine.Start()

	if err != nil {
		panic(err)
	}
}