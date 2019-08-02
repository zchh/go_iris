package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	db "iris_test/databases"
	"log"
)

//定义person类型结构
//type Person struct {
//	gorm.Model
//	Id        int    `json:"id"`
//	FirstName string `json:"first_name" validate:"required"`
//	LastName  string `json:"last_name" validate:"required"`
//}

type Person struct {
	gorm.Model
	Id        int
	FirstName string
	LastName  string
}



var dbCon = db.GetDbConnect()
var GormCon = db.GetGormConnect()

func (p *Person) AddPerson() (id int64, err error) {
	rs, err := dbCon.Exec("INSERT INTO person(first_name, last_name) VALUES (?, ?)", p.FirstName, p.LastName)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	return
}

//func (p *Person) GetPersons() (persons []Person, err error) {
//	persons = make([]Person, 0)
//	rows, err := dbCon.Query("SELECT id, first_name, last_name FROM person")
//	defer rows.Close()
//	if err != nil {
//		return
//	}
//	for rows.Next() {
//		var person Person
//		rows.Scan(&person.Id, &person.FirstName, &person.LastName)
//		persons = append(persons, person)
//	}
//	if err = rows.Err(); err != nil {
//		return
//	}
//	return
//}

func (p *Person) GetPersons(param map[string]string) (result map[int]map[string]string, err error) {
	fields,paramErr := param["fields"]
	sql := "SELECT *"
	if paramErr == true{
		sql = "SELECT "+fields
	}
	sql += " from person"
	keywords,paramErr := param["keywords"]
	if paramErr == true{
		sql += " where first_name like '%"+keywords+"%'"
	}
	rows, err := dbCon.Query(sql)
	defer rows.Close()
	if err != nil {
		return
	}
	//读出查询出的列字段名
	cols, _ := rows.Columns()
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))
	//query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	//让每一行数据都填充到[][]byte里面
	for i := range values {
		scans[i] = &values[i]
	}
	//最后得到的map
	results := make(map[int]map[string]string)
	i := 0
	for rows.Next() {
		rows.Scan(scans...)
		row := make(map[string]string) //每行数据
		for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
			key := cols[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i ++
	}
	return results,err
}

func (p *Person) GetPersonsBind(param map[string]string) (result map[int]map[string]string, err error) {
	fields,paramErr := param["fields"]
	sql := "SELECT *"
	if paramErr == true{
		sql = "SELECT "+fields
	}
	sql += " from person"
	keywords,paramErr := param["keywords"]
	println(keywords)
	if paramErr == true{
		sql += " where first_name like ?"
	}
	stmt, err := dbCon.Prepare(sql)
    rows,_ := stmt.Query(keywords)

	defer rows.Close()
	if err != nil {
		return
	}
	//读出查询出的列字段名
	cols, _ := rows.Columns()
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))
	//query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	//让每一行数据都填充到[][]byte里面
	for i := range values {
		scans[i] = &values[i]
	}
	//最后得到的map
	results := make(map[int]map[string]string)
	i := 0
	for rows.Next() {
		rows.Scan(scans...)
		row := make(map[string]string) //每行数据
		for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
			key := cols[k]
			row[key] = string(v)
		}
		results[i] = row //装入结果集中
		i ++
	}
	return results,err
}


func (p *Person) GormSelect(param map[string]string){


	//db, err := gorm.Open("mysql", "user:password@/dbname?charset=utf8&parseTime=True&loc=Local")
	dbs,err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:12330)/test?charset=utf8")
	if err != nil {
		panic("failed to connect database")
	}
	defer dbs.Close()

	var person Person
	//row := dbs.Debug().Table("person").Where(param).Select("first_name, last_name").Row()

	//firstName,paramErr := param["first_name"]





	//row := dbs.Debug().Table("person").Where("first_name = ?", "chi").Where("last_name = ?", "zz").Select("first_name, last_name").Row() // (*sql.Row)
	//row.Scan(&person.FirstName, &person.LastName)



	dbs.Select("fist_name,last_name").Find(&person)




	//aa := dbs.Debug().Unscoped().First(&product, 1).Scan(&product) // find product with id 1



	//log.Println(row)



}


func (p *Person) GetPerson() (person Person, err error) {
	err = dbCon.QueryRow("SELECT id, first_name, last_name FROM person WHERE id=?", p.Id).Scan(
		&person.Id, &person.FirstName, &person.LastName,
	)
	return
}

func (p *Person) ModPerson() (ra int64, err error) {

	//db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:12330)/test?charset=utf8")


	stmt, err := dbCon.Prepare("UPDATE person SET first_name=?, last_name=? WHERE id=?")
	defer stmt.Close()
	if err != nil {
		return
	}
	rs, err := stmt.Exec(p.FirstName, p.LastName, p.Id)
	if err != nil {
		return
	}
	ra, err = rs.RowsAffected()
	return
}

func (p *Person) DelPerson() (ra int64, err error) {
	rs, err := dbCon.Exec("DELETE FROM person WHERE id=?", p.Id)
	if err != nil {
		log.Fatalln(err)
	}
	ra, err = rs.RowsAffected()
	return
}