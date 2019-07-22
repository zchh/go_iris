package main

import (
	"fmt"
	"github.com/kataras/iris"
	"reflect"

	//"github.com/kataras/iris/context" // <- HERE
)

func main() {
	//app := iris.Default()
	//app.Get("/ping", func(ctx iris.Context) {
	//	ctx.JSON(iris.Map{
	//		"message": "pong",
	//	})
	//})
	//// listen and serve on http://0.0.0.0:8080.
	//app.Run(iris.Addr(":8080"))



	// Creates an application with default middleware:
	// logger and recovery (crash-free) middleware.
	app := iris.Default()

	types := reflect.TypeOf(app)

	fmt.Println(types)
	//app.Get("/someGet", getting)
	//app.Post("/somePost", posting)
	//app.Put("/somePut", putting)
	//app.Delete("/someDelete", deleting)
	//app.Patch("/somePatch", patching)
	//app.Head("/someHead", head)
	//app.Options("/someOptions", options)

	//app.Run(iris.Addr(":8080"))
}