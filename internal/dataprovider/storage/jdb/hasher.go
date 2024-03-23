package jdb

import (
	"app/internal/domain/hgraber"
	"context"
)

func (db *Database) GetUnHashedPages(ctx context.Context) []hgraber.Page {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := []hgraber.Page{}

	for _, t := range db.data.Data.Books {
		for _, p := range t.Pages {
			if !p.Success {
				continue
			}

			if p.Hash == "" || p.Size < 1 {
				res = append(res, p.Super(t.ID))
			}
		}
	}

	return res
}

func (db *Database) UpdatePageHash(ctx context.Context, id int, pageNumber int, hash string, size int64) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	book, ok := db.data.Data.Books[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	for i, page := range book.Pages {
		if page.PageNumber == pageNumber {
			book.Pages[i].Hash = hash
			book.Pages[i].Size = size

			db.data.Data.Books[id] = book
			db.needSave = true

			return nil
		}
	}

	return hgraber.PageNotFoundError
}
