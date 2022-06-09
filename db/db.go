package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var Con *sql.DB

func SQLExec(SQL string, Param ...interface{}) {
	_, err := Con.Exec(SQL, Param...)
	if err != nil {
		log.Panic(err.Error())
	}
}
