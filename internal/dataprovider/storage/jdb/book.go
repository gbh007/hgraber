package jdb

import (
	"app/internal/dataprovider/storage/jdb/internal/modelV2"
	"app/internal/domain/hgraber"
	"context"
	"strings"
	"time"

	"slices"
)

func (db *Database) GetUnloadedBooks(ctx context.Context) []hgraber.Book {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := []hgraber.Book{}

	for _, t := range db.data.Data.Books {
		if !t.IsFullParsed() {
			res = append(res, t.Super())
		}
	}

	return res
}

func (db *Database) NewBook(ctx context.Context, name, URL string, loaded bool) (int, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	URL = strings.TrimSpace(URL)

	if _, found := db.urlIndex[URL]; found {
		return 0, hgraber.BookAlreadyExistsError
	}

	db.lastBookID++

	db.data.Data.Books[db.lastBookID] = modelV2.RawBook{
		ID: db.lastBookID,
		Info: modelV2.RawBookInfo{
			Created: time.Now(),
			URL:     URL,
			Name:    name,
		},
	}
	db.needSave = true
	db.urlIndex[URL] = db.lastBookID

	return db.lastBookID, nil
}

func (db *Database) GetBook(ctx context.Context, id int) (hgraber.Book, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	title, ok := db.data.Data.Books[id]
	if !ok {
		return hgraber.Book{}, hgraber.BookNotFoundError
	}

	return title.Super(), nil
}

func (db *Database) GetBooks(ctx context.Context, filter hgraber.BookFilter) []hgraber.Book {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := []hgraber.Book{}

	ids := db.getTitleIDs(filter.NewFirst)
	n := len(ids)

	limit := filter.Limit
	offset := filter.Offset

	if offset > n {
		return res
	}

	if offset < 0 {
		offset = 0
	}

	if offset+limit > n {
		limit = n - offset
	}

	for _, id := range ids[offset : offset+limit] {
		if book, ok := db.data.Data.Books[id]; ok {
			res = append(res, book.Super())
		}
	}

	return res
}

func (db *Database) getTitleIDs(reverse bool) []int {
	res := make([]int, 0, len(db.data.Data.Books))

	for id := range db.data.Data.Books {
		res = append(res, id)
	}

	slices.SortStableFunc(res, func(a, b int) int {
		if reverse {
			return b - a
		}

		return a - b
	})

	return res
}

func (db *Database) UpdateBookRate(ctx context.Context, id int, rating int) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	book, ok := db.data.Data.Books[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	book.Info.Rating = rating

	db.data.Data.Books[id] = book
	db.needSave = true

	return nil
}

func (db *Database) GetBookIDByURL(ctx context.Context, url string) (int, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	id, ok := db.urlIndex[url]
	if !ok {
		return 0, hgraber.BookNotFoundError
	}

	return id, nil
}
