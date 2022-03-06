package db

import (
	"app/system"
	"context"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var _db *sql.DB

func Connect(ctx context.Context) error {
	var err error
	_db, err = sql.Open("sqlite3", "./main.db")
	if err != nil {
		system.Error(ctx, err)
		return err
	}
	_db.SetMaxOpenConns(10)
	if _, err = _db.Exec(`PRAGMA foreign_keys = ON`); err != nil {
		system.Error(ctx, err)
		return err
	}
	if _, err = _db.Exec(schemaSQL); err != nil {
		system.Error(ctx, err)
		return err
	}
	return nil
}
