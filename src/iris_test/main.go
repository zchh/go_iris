package main

import (
	"github.com/kataras/iris"

	_ "iris_test/apis"
	. "iris_test/router"
)
import "iris_test/service"

//var dbCon = db.GetDbConnect()

var sessionMgr *service.SessionMgr = nil //session管理器

func main() {
	//当整个程序完成之后关闭数据库连接
	//defer dbCon.Close()
	app := InitRouter()

	//创建session管理器,”TestCookieName”是浏览器中cookie的名字，3600是浏览器cookie的有效时间（秒）
	sessionMgr = service.NewSessionMgr("TestCookieName", 3600)

	app.Run(iris.Addr(":8080"))
}