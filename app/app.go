package app

import (
	"github.com/kataras/iris"
	"github.com/wangff15386/supermarket-go/app/good"
)

// Start the app service
func Start() {
	// Creates an application with default middleware:
	// logger and recovery (crash-free) middleware.
	app := iris.Default()

	app.Get("/good/{barcode}", good.Get)
	app.Get("/goods", good.GetAll)
	// app.Put("/somePut", putting)
	// app.Delete("/someDelete", deleting)
	// app.Patch("/somePatch", patching)
	// app.Head("/someHead", head)
	// app.Options("/someOptions", options)

	app.Run(iris.Addr(":8080"))
}
