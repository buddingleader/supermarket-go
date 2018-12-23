package main

import (
	"github.com/kataras/iris"
	"github.com/wangff15386/supermarket-go/conf"
)

func main() {
	app := iris.Default()
	app.Get("/ping", func(ctx iris.Context) {
		ctx.JSON(iris.Map{
			"message": "pong",
		})
	})
	// listen and serve on http://0.0.0.0:8080.
	// app.Run(iris.Addr(":8080"))
	app.Run(iris.Addr(conf.Config.HTTPURL))
}

// run example.go and visit http://localhost:8080/ping on browser
