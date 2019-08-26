package apis

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"github.com/kataras/iris"
	"gopkg.in/go-playground/validator.v9"
	. "iris_test/models"
)

//添加
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

//获取多条
func GetPersonsApi(c iris.Context) {
	param := make(map[string]string)
	param["page"] =  c.URLParam("page")
	param["pageSize"] =  c.URLParam("pageSize")
	param["firstName"] =  c.URLParam("firstName")
	param["lastName"] =  c.URLParam("lastName")
	param["likeUserName"] =  c.URLParam("likeUserName")
	param["userId"] =  c.URLParam("userId")
	param["fields"] = "user_name,username,department,first_name,last_name"

	var person Person
	results := person.GormJoinSelect(param)
    total := person.GormSelectCount(param)

	page,_ := strconv.ParseFloat(param["page"], 64)
	pageSize,_ := strconv.ParseFloat(param["pageSize"], 64)
	total2 := float64(total)
	lastPage :=  math.Ceil(total2/pageSize)

	c.JSON(iris.Map{
		"status":  http.StatusOK,
		"total": total,
		"per_page": pageSize,
		"current_page":page,
		"last_page":lastPage,
		"data": results,
	})
}

//获取单条
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

//修改
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

//删除
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