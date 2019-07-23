package apis

import (
	"fmt"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
	"log"
	//"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"time"

	//"github.com/gin-gonic/gin"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"gopkg.in/go-playground/validator.v9"
	. "iris_test/models"
)

var (
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID,Expires: 45*time.Minute})
)

func Login(c iris.Context)  {
	session := sess.Start(c)

	session.Set("authenticated", true)
	session.Set("user_id", 2)

	//更新过期日期与新日期
	sess.ShiftExpiration(c)

	c.JSON(iris.Map{
		"status":  http.StatusOK,
	})
}

func Logout(c iris.Context)  {

	session := sess.Start(c)
	userId,_ := session.GetInt("user_id")
    authen,_ := session.GetBoolean("authenticated")
    midd := session.GetString("midd")
	//destroy，删除整个会话数据和cookie
	//sess.Destroy(c)


	c.JSON(iris.Map{
		"status":  http.StatusOK,
		"user_id": userId,
		"authen": authen,
		"midd": midd,
		//"user_id_2": userId2,
	})
}

func CheckLogin(c iris.Context) {
	session := sess.Start(c)
	_,err := session.GetInt("user_id")
	if err != nil{
		c.JSON(iris.Map{
			"status":  http.StatusUnauthorized,
			//"user_id_2": userId2,
		})
	}else{
		c.Next()
	}
}

func Export(c iris.Context) {

	var lorem = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin lobortis, lectus dictum feugiat tempus, sem neque finibus enim, sed eleifend sem nunc ac diam. Vestibulum tempus sagittis elementum`


	doc := document.New()

	// Force the TOC to update upon opening the document
	doc.Settings.SetUpdateFieldsOnOpen(true)

	// Add a TOC
	doc.AddParagraph().AddRun().AddField(document.FieldTOC)
	// followed by a page break
	doc.AddParagraph().Properties().AddSection(wml.ST_SectionMarkNextPage)

	nd := doc.Numbering.AddDefinition()
	for i := 0; i < 9; i++ {
		lvl := nd.AddLevel()
		lvl.SetFormat(wml.ST_NumberFormatNone)
		lvl.SetAlignment(wml.ST_JcLeft)
		if i%2 == 0 {
			lvl.SetFormat(wml.ST_NumberFormatBullet)
			lvl.RunProperties().SetFontFamily("Symbol")
			lvl.SetText("")
		}
		lvl.Properties().SetLeftIndent(0.5 * measurement.Distance(i) * measurement.Inch)
	}

	// and finally paragraphs at different heading levels
	for i := 0; i < 4; i++ {
		para := doc.AddParagraph()
		para.SetNumberingDefinition(nd)
		para.Properties().SetHeadingLevel(1)
		para.AddRun().AddText("First Level")

		doc.AddParagraph().AddRun().AddText(lorem)
		for i := 0; i < 3; i++ {
			para := doc.AddParagraph()
			para.SetNumberingDefinition(nd)
			para.Properties().SetHeadingLevel(2)
			para.AddRun().AddText("Second Level")
			doc.AddParagraph().AddRun().AddText(lorem)

			para = doc.AddParagraph()
			para.SetNumberingDefinition(nd)
			para.Properties().SetHeadingLevel(3)
			para.AddRun().AddText("Third Level")
			doc.AddParagraph().AddRun().AddText(lorem)
		}
	}
	doc.SaveToFile("toc.docx")

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