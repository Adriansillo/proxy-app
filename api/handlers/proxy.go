package handlers

import (
	"encoding/json"

	"github.com/Adriansillo/proxy-app/api/middlewares"
	"github.com/kataras/iris"
)

func HandlerRedirection(app *iris.Application) {
	app.Get("/", middlewares.ProxyMiddleware, proxyHandler)
}

func proxyHandler(c iris.Context) {
	response, err := json.Marshal(middlewares.GetQue())
	if err != nil {
		c.JSON(iris.Map{"status": 400, "result": "parse error"})
		return
	}
	c.JSON(iris.Map{"result": string(response)})
}
