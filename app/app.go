package app

import "github.com/kataras/iris"

func main() {
	// Creates an application with default middleware:
	// logger and recovery (crash-free) middleware.
	app := iris.Default()

	// app.Get("/someGet", getting)
	// app.Post("/somePost", posting)
	// app.Put("/somePut", putting)
	// app.Delete("/someDelete", deleting)
	// app.Patch("/somePatch", patching)
	// app.Head("/someHead", head)
	// app.Options("/someOptions", options)

	app.Run(iris.Addr(":8080"))
}
