package jdb

import (
	"app/storage/jdb/internal/model"
	"app/system"
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"
)

type DatabaseData struct {
	Titles map[int]model.RawTitle `json:"titles"`
}

type Database struct {
	data        DatabaseData
	lastTitleID int
	uniqueURLs  map[string]struct{}
	mutex       *sync.RWMutex
	ctx         context.Context
	needSave    bool
	filename    string
}

func Init(ctx context.Context, filename string) *Database {
	return &Database{
		mutex: &sync.RWMutex{},
		data: DatabaseData{
			Titles: make(map[int]model.RawTitle),
		},
		ctx:        system.NewSystemContext(ctx, "JBD"),
		uniqueURLs: make(map[string]struct{}),
		filename:   filename,
	}
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

	newData := DatabaseData{Titles: make(map[int]model.RawTitle)}

	err = decoder.Decode(&newData)
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	db.lastTitleID = 0
	db.uniqueURLs = make(map[string]struct{})

	for id, title := range newData.Titles {
		db.uniqueURLs[strings.TrimSpace(title.URL)] = struct{}{}

		if id > db.lastTitleID {
			db.lastTitleID = id
		}
	}

	db.data = newData

	return nil
}

func (db *Database) Save(ctx context.Context, path string, force bool) error {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "jdb.Save")()

	if !db.needSave && !force {
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
