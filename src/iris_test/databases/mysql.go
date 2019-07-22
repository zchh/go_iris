package databases

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)
//因为我们需要在其他地方使用SqlDB这个变量，所以需要大写代表public
//var SqlDB *sql.DB
var Con *sql.DB

//初始化方法
func init() {
	var err error
	Con, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:12330)/test?charset=utf8")

	//SqlDB, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:12330)/test?charset=utf8")
	if err != nil {
		log.Fatal(err.Error())
	}
	//连接检测
	//err = SqlDB.Ping()
	err = Con.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	return
}

func GetDbConnect() *sql.DB {
	var err error
	Con, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:12330)/test?charset=utf8")

	//SqlDB, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:12330)/test?charset=utf8")
	if err != nil {
		log.Fatal(err.Error())
	}
	//连接检测
	//err = SqlDB.Ping()
	err = Con.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	return Con
}