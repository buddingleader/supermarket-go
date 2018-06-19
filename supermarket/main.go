package main

import (
	"supermarket-go/bussiness"
	_ "supermarket-go/routers"
)

func main() {
	// if beego.BConfig.RunMode == "dev" {
	// 	beego.BConfig.WebConfig.DirectoryIndex = true
	// 	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	// }
	// beego.Run()

	bussiness.OpenBussiness()
}
