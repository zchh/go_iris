package main

import (
	"github.com/kataras/iris"

	_ "iris_test/apis"
	. "iris_test/router"
	db "iris_test/databases"

)
import "iris_test/service"

var dbCon = db.GetDbConnect()        //原生
var GormCon = db.GetGormConnect()    //gorm

var sessionMgr *service.SessionMgr = nil //session管理器

func main() {
	//当整个程序完成之后关闭数据库连接
	defer dbCon.Close()
	defer GormCon.Close()
	app := InitRouter()

	//创建session管理器,”TestCookieName”是浏览器中cookie的名字，3600是浏览器cookie的有效时间（秒）
	sessionMgr = service.NewSessionMgr("TestCookieName", 3600)

	app.Run(iris.Addr(":8080"))
}