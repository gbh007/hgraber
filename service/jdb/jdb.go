package jdb

import (
	"app/system"
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DatabaseData struct {
	Titles map[int]Title `json:"titles"`
}

type Database struct {
	data        DatabaseData
	lastTitleID int
	mutex       *sync.RWMutex
	ctx         context.Context
	needSave    bool
}

var (
	_db *Database
)

func Init(ctx context.Context) {
	_db = &Database{
		mutex: &sync.RWMutex{},
		data: DatabaseData{
			Titles: make(map[int]Title),
		},
		ctx: system.NewSystemContext(ctx, "JBD"),
	}
}

func Get() *Database {
	return _db
}

func (db *Database) Load(ctx context.Context, path string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "jdb.Load")()

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			system.Debug(ctx, "Файл базы данных отсутствует")
			return nil
		}
		system.Error(ctx, err)
		return err
	}
	defer system.IfErrFunc(ctx, file.Close)

	decoder := json.NewDecoder(file)

	newData := DatabaseData{Titles: make(map[int]Title)}

	err = decoder.Decode(&newData)
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	db.lastTitleID = 0

	for id := range newData.Titles {
		if id > db.lastTitleID {
			db.lastTitleID = id
		}
	}

	db.data = newData

	return nil
}

func (db *Database) Save(ctx context.Context, path string) error {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "jdb.Save")()

	if !db.needSave {
		system.Debug(ctx, "Сохранение данных не требуется, пропускаю")
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		system.Error(ctx, err)
		return err
	}
	defer system.IfErrFunc(ctx, file.Close)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(db.data)
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	db.needSave = false

	return nil
}
