package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var dsn = "root:root@tcp(127.0.0.1:3306)/linkchain?charset=utf8&parseTime=true"

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return db
}
