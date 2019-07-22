package router

import (
	"github.com/kataras/iris"
    . "iris_test/apis"
	)
func InitRouter() *iris.Application {

	app := iris.Default()
	//app.Get("/", IndexApi)
	v1 := app.Party("/v1")
	{
		v1.Post("/login2", Login2)
		v1.Post("/logout2", Logout2)
		v1.Post("/person", AddPersonApi)
		v1.Get("/persons", GetPersonsApi)
		v1.Get("/person/{id:uint64}", GetPersonApi)
		v1.Put("/person/:id", ModPersonApi)
		v1.Delete("/person/{id}", DelPersonApi)
	}
	return app
}