package jdb

import (
	"app/internal/dataprovider/storage/jdb/internal/modelV2"
	"app/internal/domain/hgraber"
	"context"
)

func (db *Database) UpdateBookName(ctx context.Context, id int, name string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	book, ok := db.data.Data.Books[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	book.Info.Name = name

	db.data.Data.Books[id] = book
	db.needSave = true

	return nil

}

func (db *Database) UpdateAttributes(ctx context.Context, id int, attr hgraber.Attribute, data []string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	book, ok := db.data.Data.Books[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	values := make([]string, len(data))
	copy(values, data)

	book.Attributes[string(attr)] = modelV2.RawAttribute{
		Parsed: true,
		Values: values,
	}

	db.data.Data.Books[id] = book
	db.needSave = true

	return nil

}
