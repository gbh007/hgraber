package db

import (
	"database/sql"
	"log"
	"time"

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
	if _, err = _db.Exec(schemaSQL); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

type Title struct {
	ID               int
	Name             string
	URL              string
	PageCount        int
	CreationTime     time.Time
	Loaded           bool
	ParsedPages      bool
	ParsedTags       bool
	ParsedAuthors    bool
	ParsedCharacters bool
}

func InsertTitle(t Title) (int, error) {
	result, err := _db.Exec(
		`INSERT INTO titles(name, url, page_count, creation_time, loaded) VALUES(?, ?, ?, ?, ?)`,
		t.Name, t.URL, t.PageCount, t.CreationTime, t.Loaded,
	)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return int(id), nil
}
