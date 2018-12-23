package main

import (
	"fmt"

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

	app.Get("/profile/{name:alphabetical max(255)}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		// len(name) <=255 otherwise this route will fire 404 Not Found
		// and this handler will not be executed at all.
		fmt.Println("name:", name)
	})

	// This handler will match /user/john but will not match neither /user/ or /user.
	app.Get("/user/{name}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		ctx.Writef("Hello %s", name)
	})

	// This handler will match /users/42
	// but will not match /users/-1 because uint should be bigger than zero
	// neither /users or /users/.
	app.Get("/users/{id:uint64}", func(ctx iris.Context) {
		id := ctx.Params().GetUint64Default("id", 0)
		ctx.Writef("User with ID: %d", id)
	})

	// However, this one will match /user/john/send and also /user/john/everything/else/here
	// but will not match /user/john neither /user/john/.
	app.Post("/user/{name:string}/{action:path}", func(ctx iris.Context) {
		name := ctx.Params().Get("name")
		action := ctx.Params().Get("action")
		message := name + " is " + action
		ctx.WriteString(message)
	})

	// Query string parameters are parsed using the existing underlying request object.
	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe.
	app.Get("/welcome", func(ctx iris.Context) {
		firstname := ctx.URLParamDefault("firstname", "Guest")
		// shortcut for ctx.Request().URL.Query().Get("lastname").
		lastname := ctx.URLParam("lastname")

		ctx.Writef("Hello %s %s", firstname, lastname)
	})

	app.Post("/form_post", func(ctx iris.Context) {
		message := ctx.FormValue("message")
		nick := ctx.FormValueDefault("nick", "anonymous")

		ctx.JSON(iris.Map{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})

	app.Post("/post", func(ctx iris.Context) {
		id := ctx.URLParam("id")
		page := ctx.URLParamDefault("page", "0")
		name := ctx.FormValue("name")
		message := ctx.FormValue("message")
		// or `ctx.PostValue` for POST, PUT & PATCH-only HTTP Methods.

		app.Logger().Infof("id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})

	// listen and serve on http://0.0.0.0:8080.
	// app.Run(iris.Addr(":8080"))
	app.Run(iris.Addr(conf.Config.HTTPURL))
}

// run example.go and visit http://localhost:8080/ping on browser
