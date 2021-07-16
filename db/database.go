package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var _db *sql.DB

func Connect() error {
	var err error
	_db, err = sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Println(err)
		return err
	}
	_db.SetMaxOpenConns(10)
	if _, err = _db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		log.Println(err)
		return err
	}
	if _, err = _db.Exec(schemaSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
