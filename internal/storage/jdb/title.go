package jdb

import (
	"app/internal/domain"
	"app/internal/storage/jdb/internal/model"
	"context"
	"strings"
	"time"
)

func (db *Database) GetUnloadedTitles(ctx context.Context) []domain.Title {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := []domain.Title{}

	for _, t := range db.data.Titles {
		if !t.Data.Parsed.IsFullParsed(ctx) {
			res = append(res, t.Super(ctx))
		}
	}

	return res
}

func (db *Database) NewTitle(ctx context.Context, name, URL string, loaded bool) (int, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	URL = strings.TrimSpace(URL)

	if _, found := db.uniqueURLs[URL]; found {
		return 0, domain.TitleAlreadyExistsError
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
		return domain.TitleNotFoundError
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
		return domain.TitleNotFoundError
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

func (db *Database) GetTitle(ctx context.Context, id int) (domain.Title, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return domain.Title{}, domain.TitleNotFoundError
	}

	return title.Super(ctx), nil
}

func (db *Database) GetTitles(ctx context.Context, offset, limit int) []domain.Title {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := []domain.Title{}

	for i := db.lastTitleID - offset; i > db.lastTitleID-offset-limit; i-- {

		if title, ok := db.data.Titles[i]; ok {
			res = append(res, title.Super(ctx))
		}
	}

	return res
}

func (db *Database) GetUnsuccessedPages(ctx context.Context) []domain.PageFullInfo {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	res := []domain.PageFullInfo{}

	for _, t := range db.data.Titles {
		for i, p := range t.Pages {
			if !p.Success {
				res = append(res, domain.PageFullInfo{
					TitleID:    t.ID,
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

func (db *Database) UpdateTitlePages(ctx context.Context, id int, pages []domain.Page) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return domain.TitleNotFoundError
	}

	title.Pages = model.RawPagesFromSuper(pages)
	title.Data.Parsed.Page = len(pages) > 0

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateTitleRate(ctx context.Context, id int, rate int) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return domain.TitleNotFoundError
	}

	title.Data.Rate = rate

	db.data.Titles[id] = title
	db.needSave = true

	return nil
}

func (db *Database) GetPage(ctx context.Context, id, page int) (*domain.PageFullInfo, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return nil, domain.TitleNotFoundError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return nil, domain.PageNotFoundError
	}

	p := title.Pages[page]

	return &domain.PageFullInfo{
		TitleID:    title.ID,
		PageNumber: page + 1,
		URL:        p.URL,
		Ext:        p.Ext,
		Success:    p.Success,
		LoadedAt:   p.LoadedAt,
		Rate:       p.Rate,
	}, nil
}