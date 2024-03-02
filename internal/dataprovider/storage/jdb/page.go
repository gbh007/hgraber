package jdb

import (
	"app/internal/dataprovider/storage/jdb/internal/modelV2"
	"app/internal/domain/hgraber"
	"context"
	"time"
)

func (db *Database) UpdatePageSuccess(ctx context.Context, id, pageNumber int, success bool) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	book, ok := db.data.Data.Books[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	for i, page := range book.Pages {
		if page.PageNumber == pageNumber {
			book.Pages[i].Success = success
			if success {
				book.Pages[i].LoadAt = time.Now()
			}

			db.data.Data.Books[id] = book
			db.needSave = true

			return nil
		}
	}

	return hgraber.PageNotFoundError
}

func (db *Database) UpdatePage(ctx context.Context, id int, pageNumber int, success bool, url string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	book, ok := db.data.Data.Books[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	for i, page := range book.Pages {
		if page.PageNumber == pageNumber {
			book.Pages[i].Success = success
			book.Pages[i].URL = url
			if success {
				book.Pages[i].LoadAt = time.Now()
			}

			db.data.Data.Books[id] = book
			db.needSave = true

			return nil
		}
	}

	return hgraber.PageNotFoundError
}

func (db *Database) UpdatePageRate(ctx context.Context, id, pageNumber int, rating int) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	book, ok := db.data.Data.Books[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	for i, page := range book.Pages {
		if page.PageNumber == pageNumber {
			book.Pages[i].Rating = rating
			db.data.Data.Books[id] = book
			db.needSave = true

			return nil
		}
	}

	return hgraber.PageNotFoundError
}

func (db *Database) GetUnsuccessPages(ctx context.Context) []hgraber.Page {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := []hgraber.Page{}

	for _, t := range db.data.Data.Books {
		for i, p := range t.Pages {
			if !p.Success {
				res = append(res, hgraber.Page{
					BookID:     t.ID,
					PageNumber: i + 1,
					URL:        p.URL,
					Ext:        p.Ext,
					Success:    p.Success,
					LoadedAt:   p.LoadAt,
				})
			}
		}
	}

	return res
}

func (db *Database) UpdateBookPages(ctx context.Context, id int, pages []hgraber.Page) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	book, ok := db.data.Data.Books[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	book.Pages = modelV2.RawPagesFromSuper(pages)
	book.Info.PageCount = len(pages)

	db.data.Data.Books[id] = book
	db.needSave = true

	return nil
}

func (db *Database) GetPage(ctx context.Context, id, pageNumber int) (*hgraber.Page, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	title, ok := db.data.Data.Books[id]
	if !ok {
		return nil, hgraber.BookNotFoundError
	}

	for _, page := range title.Pages {
		if page.PageNumber == pageNumber {
			domainPage := page.Super(id)
			return &domainPage, nil
		}
	}

	return nil, hgraber.PageNotFoundError
}
