package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	db "iris_test/databases"
	"log"
	"strconv"
)

type Person struct {
	Id        int
	FirstName string
	LastName  string
}

var dbCon = db.GetDbConnect()        //原生
var GormCon = db.GetGormConnect()    //gorm

/**
    添加
    coder: 张池
 */
func (p *Person) AddPerson() (id int64, err error) {
	rs, err := dbCon.Exec("INSERT INTO person(first_name, last_name) VALUES (?, ?)", p.FirstName, p.LastName)
	if err != nil {
		return
	}
	id, err = rs.LastInsertId()
	return
}

/**
   按条件搜索人
   coder: 张池
 */
func (p *Person) GormJoinSelect(param map[string]string) map[int]map[string]string{
	db := GormCon
	defer db.Close()
	selects := db.Table("person").Joins("left join userinfo on person.user_id = userinfo.uid")
	fields,paramErr := param["fields"]
	if paramErr == true{
		selects = selects.Select(fields)
	}else{
		selects = selects.Select("*")
	}
	firstName,paramErr := param["firstName"]
	if paramErr == true && firstName != ""{
		selects = selects.Where("first_name = ?", firstName)
	}
	lastName,paramErr := param["lastName"]
	if paramErr == true && lastName != ""{
		selects = selects.Where("last_name = ?", lastName)
	}
	userName,paramErr := param["likeUserName"]
	if paramErr == true && userName != ""{
		userName = "%"+userName+"%"
		selects = selects.Where("user_name like ?", userName)
	}
	userId,paramErr := param["userId"]
	if paramErr == true && userId != ""{
		selects = selects.Where("user_id = ?", userId)
	}

	page,pageErr := param["page"]
	pageSize,sizeErr := param["pageSize"]
	if pageErr == true && sizeErr == true {
		page,_ := strconv.Atoi(page)
		pageSize,_ := strconv.Atoi(pageSize)
		offset := (page-1)*pageSize

		selects = selects.Limit(pageSize).Offset(offset)
	}

	rows,_:= selects.Order("id asc").Rows()

	//读出查询出的列字段名
	cols, _ := rows.Columns()
	//query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
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
	return results
}

func (p *Person) GormSelect(param map[string]string) map[int]map[string]string{
	db,err := gorm.Open("mysql", "root:123456@tcp(127.0.0.1:12330)/test?charset=utf8")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	selects := db.Table("person")
	fields,paramErr := param["fields"]
	if paramErr == true{
		selects = selects.Select(fields)
	}else{
		selects = selects.Select("*")
	}
	firstName,paramErr := param["firstName"]
	if paramErr == true && firstName != ""{
		selects = selects.Where("first_name = ?", firstName)
	}
	lastName,paramErr := param["lastName"]
	if paramErr == true && lastName != ""{
		selects = selects.Where("last_name = ?", lastName)
	}
	userName,paramErr := param["likeUserName"]
	if paramErr == true && userName != ""{
		userName = "%"+userName+"%"
		selects = selects.Where("user_name like ?", userName)
	}

	page,pageErr := param["page"]
	pageSize,sizeErr := param["pageSize"]
    if pageErr == true && sizeErr == true {
		page,_ := strconv.Atoi(page)
		pageSize,_ := strconv.Atoi(pageSize)
		offset := (page-1)*pageSize

		selects = selects.Limit(pageSize).Offset(offset)
	}

	rows,err:= selects.Order("id asc").Rows()

	//读出查询出的列字段名
	cols, _ := rows.Columns()
	//query.Scan的参数，因为每次查询出来的列是不定长的，用len(cols)定住当次查询的长度
	scans := make([]interface{}, len(cols))
	//values是每个列的值，这里获取到byte里
	values := make([][]byte, len(cols))
	for i := range values {
		scans[i] = &values[i]
	}
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
	return results
}

/**
   按条件计算搜索的数量
   coder: 张池
 */
func (p *Person) GormSelectCount(param map[string]string) int{
	db := GormCon
	defer db.Close()
	selects := db.Table("person")
	firstName,paramErr := param["firstName"]
	if paramErr == true && firstName != ""{
		selects = selects.Where("first_name = ?", firstName)
	}
	lastName,paramErr := param["lastName"]
	if paramErr == true && lastName != ""{
		selects = selects.Where("last_name = ?", lastName)
	}
	userName,paramErr := param["like_user_name"]
	if paramErr == true && userName != ""{
		userName = "%"+userName+"%"
		selects = selects.Where("user_name like ?", userName)
	}
	var count int
	selects.Count(&count)
	return count
}



func (p *Person) GetPerson() (person Person, err error) {
	err = dbCon.QueryRow("SELECT id, first_name, last_name FROM person WHERE id=?", p.Id).Scan(
		&person.Id, &person.FirstName, &person.LastName,
	)
	return
}

/**
   更新
   coder: 张池
 */
func (p *Person) ModPerson() (ra int64, err error) {
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

/**
   删除
   coder: 张池
 */
func (p *Person) DelPerson() (ra int64, err error) {
	rs, err := dbCon.Exec("DELETE FROM person WHERE id=?", p.Id)
	if err != nil {
		log.Fatalln(err)
	}
	ra, err = rs.RowsAffected()
	return
}