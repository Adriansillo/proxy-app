package server

import (
	"os"

	"github.com/kataras/iris"
)

func SetUp() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("debug")
	return app
}

func RunServer(app *iris.Application) {
	app.Run(iris.Addr(os.Getenv("PORT")))
}
