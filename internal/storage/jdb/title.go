package jdb

import (
	"app/internal/domain"
	"app/internal/storage/jdb/internal/model"
	"context"
	"strings"
	"time"

	"slices"
)

func (db *Database) GetUnloadedBooks(ctx context.Context) []domain.Book {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := []domain.Book{}

	for _, t := range db.data.Titles {
		if !t.Data.Parsed.IsFullParsed() {
			res = append(res, t.Super())
		}
	}

	return res
}

func (db *Database) NewBook(ctx context.Context, name, URL string, loaded bool) (int, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	URL = strings.TrimSpace(URL)

	if _, found := db.uniqueURLs[URL]; found {
		return 0, domain.BookAlreadyExistsError
	}

	db.lastTitleID++

	db.data.Titles[db.lastTitleID] = model.RawTitle{
		ID:      db.lastTitleID,
		Created: time.Now(),
		URL:     URL,
		Pages:   []model.RawPage{},
		Data: model.RawTitleInfo{
			Parsed: model.RawTitleInfoParsed{
				Name: loaded,
			},
			Name:       name,
			Tags:       []string{},
			Authors:    []string{},
			Characters: []string{},
			Languages:  []string{},
			Categories: []string{},
			Parodies:   []string{},
			Groups:     []string{},
		},
	}
	db.needSave = true
	db.uniqueURLs[URL] = struct{}{}

	return db.lastTitleID, nil
}

func (db *Database) UpdatePageSuccess(ctx context.Context, id, page int, success bool) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return domain.BookNotFoundError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return domain.PageNotFoundError
	}

	title.Pages[page].Success = success
	if success {
		title.Pages[page].LoadedAt = time.Now()
	}

	db.data.Titles[id] = title
	db.needSave = true

	return nil
}

func (db *Database) UpdatePageRate(ctx context.Context, id, page int, rate int) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return domain.BookNotFoundError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return domain.PageNotFoundError
	}

	title.Pages[page].Rate = rate

	db.data.Titles[id] = title
	db.needSave = true

	return nil
}

func (db *Database) GetBook(ctx context.Context, id int) (domain.Book, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return domain.Book{}, domain.BookNotFoundError
	}

	return title.Super(), nil
}

func (db *Database) GetBooks(ctx context.Context, filter domain.BookFilter) []domain.Book {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := []domain.Book{}

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
		if title, ok := db.data.Titles[id]; ok {
			res = append(res, title.Super())
		}
	}

	return res
}

func (db *Database) getTitleIDs(reverse bool) []int {
	res := make([]int, 0, len(db.data.Titles))

	for id := range db.data.Titles {
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

func (db *Database) GetUnsuccessPages(ctx context.Context) []domain.Page {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := []domain.Page{}

	for _, t := range db.data.Titles {
		for i, p := range t.Pages {
			if !p.Success {
				res = append(res, domain.Page{
					BookID:     t.ID,
					PageNumber: i + 1,
					URL:        p.URL,
					Ext:        p.Ext,
					Success:    p.Success,
					LoadedAt:   p.LoadedAt,
				})
			}
		}
	}

	return res
}

func (db *Database) UpdateBookPages(ctx context.Context, id int, pages []domain.Page) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return domain.BookNotFoundError
	}

	title.Pages = model.RawPagesFromSuper(pages)
	title.Data.Parsed.Page = len(pages) > 0

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateBookRate(ctx context.Context, id int, rate int) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return domain.BookNotFoundError
	}

	title.Data.Rate = rate

	db.data.Titles[id] = title
	db.needSave = true

	return nil
}

func (db *Database) GetPage(ctx context.Context, id, page int) (*domain.Page, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return nil, domain.BookNotFoundError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return nil, domain.PageNotFoundError
	}

	p := title.Pages[page]

	return &domain.Page{
		BookID:     title.ID,
		PageNumber: page + 1,
		URL:        p.URL,
		Ext:        p.Ext,
		Success:    p.Success,
		LoadedAt:   p.LoadedAt,
		Rate:       p.Rate,
	}, nil
}
