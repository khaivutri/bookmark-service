package main

import (
	"github.com/khaivutri/bookmark-service/infrastructure"
)

//	@title Bookmark Service API
//	@version 2.0
//	@description This is a simple REST API for a bookmark service - Demo MLIoT Lab.
//	@BasePath /
//	@securityDefinitions.apikey BearerAuth
//	@in header
//	@name Authorization
func main() {
	app := infrastructure.CreateApp()
	err := app.Start()
	if err != nil {
		panic(err)
	}
}

