package jdb

import (
	"app/internal/dataprovider/storage/jdb/internal/model"
	"app/pkg/ctxtool"
	"app/pkg/logger"
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
	filename    *string

	logger *logger.Logger
}

func Init(ctx context.Context, logger *logger.Logger, filename *string) *Database {
	return &Database{
		mutex: &sync.RWMutex{},
		data: DatabaseData{
			Titles: make(map[int]model.RawTitle),
		},
		ctx:        ctxtool.NewSystemContext(ctx, "JBD"),
		uniqueURLs: make(map[string]struct{}),
		filename:   filename,
		logger:     logger,
	}
}

func (db *Database) Load(ctx context.Context, path string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			db.logger.Debug(ctx, "Файл базы данных отсутствует")
			return nil
		}
		db.logger.Error(ctx, err)
		return err
	}
	defer db.logger.IfErrFunc(ctx, file.Close)

	decoder := json.NewDecoder(file)

	newData := DatabaseData{Titles: make(map[int]model.RawTitle)}

	err = decoder.Decode(&newData)
	if err != nil {
		db.logger.Error(ctx, err)
		return err
	}

	db.lastTitleID = 0
	db.uniqueURLs = make(map[string]struct{})

	for id, title := range newData.Titles {
		u := strings.TrimSpace(title.URL)
		if _, found := db.uniqueURLs[u]; found {
			db.logger.Warning(ctx, "Дублирование ссылки при загрузке БД", u)
		} else {
			db.uniqueURLs[u] = struct{}{}
		}

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

	if !db.needSave && !force {
		db.logger.Debug(ctx, "Сохранение данных не требуется, пропускаю")

		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		db.logger.Error(ctx, err)

		return err
	}

	defer db.logger.IfErrFunc(ctx, file.Close)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(db.data)
	if err != nil {
		db.logger.Error(ctx, err)
		return err
	}

	db.needSave = false

	return nil
}
