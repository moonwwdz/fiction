package db

import (
	"database/sql"
	"time"

	"github.com/moonwwdz/fiction/util"

	_ "github.com/mattn/go-sqlite3"
)

var tableName string
var conn *sql.DB

func init() {
	var err error
	conn, err = sql.Open("sqlite3", "./foo.db")
	if err != nil {
		panic(err)
	}
}

type ficS struct {
	Uid        int
	Title      string
	TitleMd5   string
	Cont       string
	CreateTime *time.Time
}

//SetTable 设置表名的同时，创建表
func SetTable(sqlTableName string) {
	tableName = sqlTableName
	_, err := conn.Query("select * from " + tableName)
	if err == nil {
		return
	}

	tableSql := "CREATE TABLE IF NOT EXISTS " + tableName + "("
	tableSql += "uid INTEGER PRIMARY KEY AUTOINCREMENT,"
	tableSql += "title VARCHAR(64) NULL,"
	tableSql += "titleMd5 varchar(64),"
	tableSql += "cont text NULL,"
	tableSql += "createTime TimeStamp NOT NULL DEFAULT (datetime('now','localtime')) );"

	conn.Exec(tableSql)
}

//Savecont 保存内容
func SaveCont(title, cont string) (int64, error) {
	preSql := "insert into " + tableName + "(title,titleMd5,cont) values (?,?,?)"
	stmt, err := conn.Prepare(preSql)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(title, util.GetMD5Hash(title), cont)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return id, nil
}

//Getcontbymd5 取内容
func GetContByMd5(md5Str string) (string, error) {
	rows, err := conn.Query("select cont from " + tableName + " where titleMd5 = '" + md5Str + "'")
	if err != nil {
		return "", err
	}

	var cont string
	for rows.Next() {
		err = rows.Scan(&cont)
	}
	rows.Close()
	return cont, nil
}

//Getlasterfive 取最近五条记录
func GetLasterFive() ([]ficS, error) {
	sqlStr := "select * from " + tableName + " order by createtime desc limit 5"
	rows, err := conn.Query(sqlStr)
	if err != nil {
		return nil, err
	}

	var ficList []ficS
	for rows.Next() {
		fic := new(ficS)
		err = rows.Scan(&fic.Uid, &fic.Title, &fic.TitleMd5, &fic.Cont, &fic.CreateTime)
		ficList = append(ficList, *fic)
	}
	return ficList, nil
}
