package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func initdb(dbpath string) (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", dbpath)
	return
}
