package apis

import (
	"fmt"
	"log"

	//"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"

	//"github.com/gin-gonic/gin"
	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
	. "iris_test/models"
	"iris_test/service"
    "github.com/kataras/iris/sessions"
)

var (
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
)


var sessionMgr *service.SessionMgr = nil //session管理器


func Login(c iris.Context)  {
	//session 设置
	//firstName := c.FormValue("user_name")
	//lastName := c.FormValue("password")

	var userSession = make(map[interface{}]interface{})

	userSession["user_id"] = "111"


	session := service.SessionStore{"123",time.Now(),userSession}

    session.Set("user_id", "1")

	userId := session.Get("user_id")


	c.JSON(iris.Map{
		"status":  http.StatusOK,
		"user_id": userId,
	})
}

func Logout(c iris.Context)  {
	var session service.SessionStore
	user := session.Get("user_id")
	c.JSON(iris.Map{
		"status":  http.StatusOK,
		"user": user,
	})
}

func Login2(c iris.Context)  {
	session := sess.Start(c)

	session.Set("authenticated", true)
	c.JSON(iris.Map{
		"status":  http.StatusOK,
	})
}

func Logout2(c iris.Context, w http.ResponseWriter, r *http.Request)  {

	var sessionID = sessionMgr.CheckCookieValid(w, r)

	sessionMgr.EndSession(w, r) //用户退出时删除对应session
	http.Redirect(w, r, "/login", http.StatusFound)
	c.JSON(iris.Map{
		"status":  http.StatusOK,
		"sessionID": sessionID,
	})
}

func todayFilename() string {
	today := time.Now().Format("Jan 02 2006")
	return today + ".txt"
}

func newLogFile() *os.File {
	filename := todayFilename()
	// Open the file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}


func AddPersonApi(c iris.Context) {
	var validate *validator.Validate
	validate = validator.New()
	firstName := c.FormValue("first_name")
	lastName := c.FormValue("last_name")
	person := Person{FirstName: firstName, LastName: lastName}
	err := validate.Struct(person)
	if err != nil {
		c.StatusCode(iris.StatusBadRequest)
		errMsg := ""
		for _, err := range err.(validator.ValidationErrors) {
			if len(errMsg) == 0{
				errMsg = err.StructField() + " " + err.Tag()
			}else{
				errMsg = errMsg + "," + err.StructField() + " " + err.Tag()
			}
		}
		c.JSON(iris.Map{
			"status":  http.StatusNoContent,
		    "msg": errMsg,
		})
		// from here you can create your own error messages in whatever language you wish.
		return
	}

	ra, err := person.AddPerson()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("insert successful %d", ra)

	c.JSON(iris.Map{
		"status":  http.StatusOK,
		"msg": msg,
	})
}

func GetPersonsApi(c iris.Context) {
	var p Person
	persons, err := p.GetPersons()
	if err != nil {
		log.Fatalln(err)
	}

	c.JSON(iris.Map{
		"status":  http.StatusOK,
		"persons": persons,
	})
}

func GetPersonApi(c iris.Context) {
    cid := c.Params().Get("id")

	//cid := c.URLParam("id")
	//cid := c.URLParamDefault("id","1")
	//c.Writef("Hello %s", cid)

	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	p := Person{Id: id}
	person, err := p.GetPerson()
	if err != nil {
		log.Fatalln(err)
	}

	c.JSON(iris.Map{
		"status":  http.StatusOK,
		"person": person,
	})
}

func ModPersonApi(c iris.Context) {
	cid := c.URLParam("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	p := Person{Id: id}
	//err = c.Bind(&p)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	ra, err := p.ModPerson()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("Update person %d successful %d", p.Id, ra)

	c.JSON(iris.Map{
		"status":  http.StatusOK,
		"msg": msg,
	})
}

func DelPersonApi(c iris.Context) {

	//记录日志1
	f := newLogFile()
	defer f.Close()
	app := iris.New()
	app.Logger().SetOutput(f)
	//c.WriteString("pong")

	cid := c.Params().Get("id")
	id, err := strconv.Atoi(cid)
	if err != nil {
		log.Fatalln(err)
	}
	p := Person{Id: id}
	ra, err := p.DelPerson()
	if err != nil {
		log.Fatalln(err)
	}
	msg := fmt.Sprintf("Delete person %d successful %d", id, ra)

	//记录日志2
	c.Application().Logger().Infof("delete person id = %d",id)

	c.JSON(iris.Map{
		"status":  http.StatusOK,
		"msg": msg,
	})
}