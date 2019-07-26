package apis

import (
	"fmt"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
	"iris_test/service"
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

func ExportWordByTemp(c iris.Context)  {

	var lorem = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Proin lobortis, lectus dictum feugiat tempus, sem neque finibus enim, sed eleifend sem nunc ac diam. Vestibulum tempus sagittis elementum`


	// When Word saves a document, it removes all unused styles.  This means to
	// copy the styles from an existing document, you must first create a
	// document that contains text in each style of interest.  As an example,
	// see the template.docx in this directory.  It contains a paragraph set in
	// each style that Word supports by default.
	doc, err := document.OpenTemplate("resource/template.docx")
	if err != nil {
		log.Fatalf("error opening Windows Word 2016 document: %s", err)
	}

	// We can now print out all styles in the document, verifying that they
	// exist.
	for _, s := range doc.Styles.Styles() {
		fmt.Println("style", s.Name(), "has ID of", s.StyleID(), "type is", s.Type())
	}

	// And create documents setting their style to the style ID (not style name).
	para := doc.AddParagraph()
	para.SetStyle("Title")
	para.AddRun().AddText("My Document Title")

	para = doc.AddParagraph()
	para.SetStyle("Subtitle")
	para.AddRun().AddText("Document Subtitle")

	para = doc.AddParagraph()
	para.SetStyle("Heading1")
	para.AddRun().AddText("Major Section")
	para = doc.AddParagraph()
	para = doc.AddParagraph()
	for i := 0; i < 4; i++ {
		para.AddRun().AddText(lorem)
	}

	para = doc.AddParagraph()
	para.SetStyle("Heading2")
	para.AddRun().AddText("Minor Section")
	para = doc.AddParagraph()
	for i := 0; i < 4; i++ {
		para.AddRun().AddText(lorem)
	}

	// using a pre-defined table style
	table := doc.AddTable()
	table.Properties().SetWidthPercent(90)
	table.Properties().SetStyle("GridTable4-Accent1")
	look := table.Properties().TableLook()
	// these have default values in the style, so we manually turn some of them off
	look.SetFirstColumn(false)
	look.SetFirstRow(true)
	look.SetLastColumn(false)
	look.SetLastRow(true)
	look.SetHorizontalBanding(true)

	for r := 0; r < 5; r++ {
		row := table.AddRow()
		for c := 0; c < 5; c++ {
			cell := row.AddCell()
			cell.AddParagraph().AddRun().AddText(fmt.Sprintf("row %d col %d", r+1, c+1))
		}
	}
	doc.SaveToFile("use-template.docx")


	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Disposition", "attachment; filename="+"use-template.docx")//文件名
	c.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	//最主要的一句
	http.ServeFile(c.ResponseWriter(), c.Request(),"use-template.docx")

}

func ExportPerson(c iris.Context)  {
	header := [...]string{"id","名","姓"}

	var p Person
	persons, err := p.GetPersons()
	if err != nil {
		log.Fatalln(err)
	}
	service.ExportExcel(c, header, persons)

}


func Download(c iris.Context){
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Disposition", "attachment; filename="+"toc.docx")//文件名
	c.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	//最主要的一句
	http.ServeFile(c.ResponseWriter(), c.Request(),"toc.docx")
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