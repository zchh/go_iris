package apis

import (
	"fmt"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
	"github.com/unidoc/unioffice/spreadsheet"
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

func ExportExcel(c iris.Context)  {
	ss := spreadsheet.New()
	sheet := ss.AddSheet()

	header := [...]string{"用户名","部门","创建时间"}
    cellArr := [...]string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}

	row := sheet.AddRow()   //下移
	for i := 0; i < len(header); i++{


        format := cellArr[i]+"0"

		cell := row.AddNamedCell(fmt.Sprintf(format))

		//cell := row.AddNamedCell(fmt.Sprintf("%c", 'A'+i))
		cell.SetString(fmt.Sprintf("%s", header[i]))
	}
	row2 := sheet.AddRow()   //下移

	mapss := map[string]interface{}{"username":"aaa","department":"研发部门","created":"2012-12-09"}
	mapss2 := map[string]interface{}{"username":"bbb","department":"研发部门","created":"2012-12-09"}

	mapArr2 := map[int]interface{}{0:mapss,1:mapss2}


	for _,value := range mapArr2 {

		for _,value2 := range value{


		}


		fmt.Println(value)
	}



	for j := 0;j < len(header); j++{

		for d := 1; d <= len(mapArr2); d ++ {
			format2 := cellArr[j]+string(d)

			cell := row2.AddNamedCell(fmt.Sprintf(format2))



			for _,value := range mapArr2 {
				fmt.Println(value)
			}

			//cell := row.AddNamedCell(fmt.Sprintf("%c", 'A'+i))
			cell.SetString(fmt.Sprintf("%s", header[i]))

		}



	}

	for r := 0; r < 5; r++ {
		row := sheet.AddRow()   //下移

		// can't add an un-named cell to row zero here as we also add cell 'A1',
		// meaning the un-naned cell must come before 'A1' which is invalid.
		if r != 0 {
			// an unnamed cell displays in the first available column
			row.AddCell().SetString("unnamed-before") //右移赋值
		}

		// setting these to A, B, C, specifically
		cell := row.AddNamedCell(fmt.Sprintf("%c", 'A'+r))
		cell.SetString(fmt.Sprintf("row %d", r))

		// an un-named cell after a named cell is display immediately after a named cell
		row.AddCell().SetString("unnamed-after")
	}

	sheet.AddNumberedRow(26).AddNamedCell("C").SetString("Cell C26")

	// This line would create an invalid sheet with two identically ID'd rows
	// which would fail validation below
	// sheet.AddNumberedRow(26).AddNamedCell("C27").SetString("Cell C27")

	// so instead use Row which will create or retrieve an existing row
	sheet.Row(26).AddNamedCell("E").SetString("Cell E26")
	sheet.Row(26).Cell("F").SetString("Cell F26")

	// You can also reference cells fully from the sheet.
	sheet.Cell("H1").SetString("Cell H1")
	sheet.Cell("H2").SetString("Cell H2")
	sheet.Cell("H3").SetString("Cell H3")

	if err := ss.Validate(); err != nil {
		log.Fatalf("error validating sheet: %s", err)
	}

	ss.SaveToFile("resource/named-cells.xlsx")

	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Disposition", "attachment; filename="+"named-cells.xlsx")//文件名
	c.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	//最主要的一句
	http.ServeFile(c.ResponseWriter(), c.Request(),"resource/named-cells.xlsx")
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