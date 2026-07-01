package main

import "github.com/khaivutri/bookmark-service/internal/api"

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