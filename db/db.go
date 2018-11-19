package db

import "database/sql"

var tableName string

func init() {
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}
}

//SetTable 设置表名的同时，创建表
func SetTable(tableName string) {
	tableName = tableName
}

//Savecont 保存内容
func SaveCont(title, cont string) (id int, err error) {
	stmt, err := db.Prepare("insert into " + tableName + "(title,title_md5,cont) values (?,?,?)")
	if err != nil {
		return nil, err
	}

	res, err := stmt.Exec(title, util.GetMd5Hash(title), cont)
	if err != nil {
		return nil, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}
