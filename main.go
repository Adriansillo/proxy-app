package main

import (
	"github.com/Adriansillo/proxy-app/api/handlers"
	"github.com/Adriansillo/proxy-app/api/server"
	"github.com/Adriansillo/proxy-app/api/utils"
)

func main() {
	utils.LoadEnv()
	app := server.SetUp()
	handlers.HandlerRedirection(app)
	server.RunServer(app)
}
