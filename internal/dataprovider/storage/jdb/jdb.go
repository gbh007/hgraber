package jdb

import (
	"app/internal/dataprovider/storage/jdb/internal/fileModel"
	"app/pkg/ctxtool"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type Database struct {
	data       *fileModel.DatabaseData
	lastBookID int
	urlIndex   map[string]int
	mutex      *sync.RWMutex
	ctx        context.Context
	needSave   bool
	filename   *string

	logger *slog.Logger
}

func Init(ctx context.Context, logger *slog.Logger, filename *string) *Database {
	return &Database{
		mutex:    &sync.RWMutex{},
		data:     fileModel.New(),
		ctx:      ctxtool.NewSystemContext(ctx, "JBD"),
		urlIndex: make(map[string]int),
		filename: filename,
		logger:   logger,
	}
}

func (db *Database) Load(ctx context.Context, path string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			db.logger.DebugContext(ctx, "Файл базы данных отсутствует")
			return nil
		}
		db.logger.ErrorContext(ctx, err.Error())
		return err
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			db.logger.ErrorContext(ctx, closeErr.Error())
		}
	}()

	decoder := json.NewDecoder(file)

	newData := fileModel.New()

	err = decoder.Decode(&newData)
	if err != nil {
		db.logger.ErrorContext(ctx, err.Error())
		return err
	}

	isMigrated, err := newData.Migrate()
	if err != nil {
		db.logger.ErrorContext(ctx, err.Error())
		return err
	}

	if isMigrated {
		db.logger.WarnContext(ctx, "Произведена миграция")
		db.needSave = true
	}

	db.lastBookID = 0
	db.urlIndex = make(map[string]int)

	for id, book := range newData.Data.Books {
		u := strings.TrimSpace(book.Info.URL)
		if _, found := db.urlIndex[u]; found {
			db.logger.WarnContext(ctx, "Дублирование ссылки при загрузке БД", slog.String("url", u))
		} else {
			db.urlIndex[u] = id
		}

		if id > db.lastBookID {
			db.lastBookID = id
		}
	}

	db.data = newData

	return nil
}

func (db *Database) Save(ctx context.Context, path string, force bool) error {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	if !db.needSave && !force {
		db.logger.DebugContext(ctx, "Сохранение данных не требуется, пропускаю")

		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		db.logger.ErrorContext(ctx, err.Error())

		return err
	}

	defer func() {
		closeErr := file.Close()
		if closeErr != nil {
			db.logger.ErrorContext(ctx, closeErr.Error())
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(db.data)
	if err != nil {
		db.logger.ErrorContext(ctx, err.Error())
		return err
	}

	db.needSave = false

	return nil
}
