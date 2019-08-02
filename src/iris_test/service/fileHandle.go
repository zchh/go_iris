package service

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/unidoc/unioffice/spreadsheet"
	"log"
	"net/http"
	"os"
)

//data map[int]map[string]string
func ExportExcel(c iris.Context, header []string, data map[int]map[string]string)  {
	ss := spreadsheet.New()
	sheet := ss.AddSheet()
	//header := [...]string{"用户名","部门","创建时间"}
	cellArr := [...]string{"A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z"}
	row := sheet.AddRow()   //下移
	for i := 0; i < len(header); i++{
		format := cellArr[i]+"0"
		cell := row.AddNamedCell(fmt.Sprintf(format))
		cell.SetString(fmt.Sprintf("%s", header[i]))
	}
	row = sheet.AddRow()   //下移
	//mapss := map[string]string{"username":"aaa","department":"研发部门","created":"2012-12-09"}
	//mapss2 := map[string]string{"username":"bbb","department":"研发部门","created":"2012-12-09"}
	//mapArr2 := map[int]map[string]string{0:mapss,1:mapss2}

	mapArr2 := data

	j := 0
	for _,value := range mapArr2 {
		for _,value2 := range value{
			format2 := cellArr[j]+"0"
			cell := row.AddNamedCell(fmt.Sprintf(format2))
			cell.SetString(fmt.Sprintf("%s", value2))
			j++
		}
		row = sheet.AddRow()   //下移
		j = 0
	}
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

	os.Remove("resource/named-cells.xlsx")
}
