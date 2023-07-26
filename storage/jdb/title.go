package jdb

import (
	"app/storage/jdb/internal/model"
	"app/storage/schema"
	"app/system"
	"context"
	"strings"
	"time"
)

func (db *Database) GetUnloadedTitles(ctx context.Context) []schema.Title {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "UnloadedTitles")()

	res := []schema.Title{}

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

	defer system.Stopwatch(ctx, "NewTitle")()

	URL = strings.TrimSpace(URL)

	if _, found := db.uniqueURLs[URL]; found {
		return 0, schema.TitleDuplicateError
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

	defer system.Stopwatch(ctx, "UpdatePageSuccess")()

	title, ok := db.data.Titles[id]
	if !ok {
		return schema.TitleIndexError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return schema.PageIndexError
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

	defer system.Stopwatch(ctx, "UpdatePageRate")()

	title, ok := db.data.Titles[id]
	if !ok {
		return schema.TitleIndexError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return schema.PageIndexError
	}

	title.Pages[page].Rate = rate

	db.data.Titles[id] = title
	db.needSave = true

	return nil
}

func (db *Database) GetTitle(ctx context.Context, id int) (schema.Title, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "GetTitle")()

	title, ok := db.data.Titles[id]
	if !ok {
		return schema.Title{}, schema.TitleIndexError
	}

	return title.Super(ctx), nil
}

func (db *Database) GetTitles(ctx context.Context, offset, limit int) []schema.Title {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "GetTitles")()

	res := []schema.Title{}

	for i := db.lastTitleID - offset; i > db.lastTitleID-offset-limit; i-- {

		if title, ok := db.data.Titles[i]; ok {
			res = append(res, title.Super(ctx))
		}
	}

	return res
}

func (db *Database) GetUnsuccessedPages(ctx context.Context) []schema.PageFullInfo {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "GetUnsuccessPages")()

	res := []schema.PageFullInfo{}

	for _, t := range db.data.Titles {
		for i, p := range t.Pages {
			if !p.Success {
				res = append(res, schema.PageFullInfo{
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

func (db *Database) UpdateTitlePages(ctx context.Context, id int, pages []schema.Page) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "UpdateTitlePages")()

	title, ok := db.data.Titles[id]
	if !ok {
		return schema.TitleIndexError
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

	defer system.Stopwatch(ctx, "UpdateTitleRate")()

	title, ok := db.data.Titles[id]
	if !ok {
		return schema.TitleIndexError
	}

	title.Data.Rate = rate

	db.data.Titles[id] = title
	db.needSave = true

	return nil
}

func (db *Database) GetPage(ctx context.Context, id, page int) (*schema.PageFullInfo, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "GetPage")()

	title, ok := db.data.Titles[id]
	if !ok {
		return nil, schema.TitleIndexError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return nil, schema.PageIndexError
	}

	p := title.Pages[page]

	return &schema.PageFullInfo{
		TitleID:    title.ID,
		PageNumber: page + 1,
		URL:        p.URL,
		Ext:        p.Ext,
		Success:    p.Success,
		LoadedAt:   p.LoadedAt,
		Rate:       p.Rate,
	}, nil
}
