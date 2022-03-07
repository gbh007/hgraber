package jdb

import (
	"app/system"
	"context"
	"time"
)

func (db *Database) GetUnloadedTitles(ctx context.Context) []Title {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "UnloadedTitles")()

	res := []Title{}

	for _, t := range db.data.Titles {
		if !t.Data.Parsed.IsFullParsed(ctx) {
			res = append(res, t.Copy(ctx))
		}
	}

	return res
}

func (db *Database) NewTitle(ctx context.Context, name, URL string, loaded bool) (int, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "NewTitle")()

	if _, found := db.uniqueURLs[URL]; found {
		return 0, TitleDuplicateError
	}

	db.lastTitleID++

	db.data.Titles[db.lastTitleID] = Title{
		ID:      db.lastTitleID,
		Created: time.Now(),
		URL:     URL,
		Pages:   []Page{},
		Data: TitleInfo{
			Parsed: TitleInfoParsed{
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
		return TitleIndexError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return PageIndexError
	}

	title.Pages[page].Success = success
	if success {
		title.Pages[page].LoadedAt = time.Now()
	}

	db.data.Titles[id] = title
	db.needSave = true

	return nil
}

func (db *Database) GetTitle(ctx context.Context, id int) (Title, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "GetTitle")()

	title, ok := db.data.Titles[id]
	if !ok {
		return Title{}, TitleIndexError
	}

	return title.Copy(ctx), nil
}

func (db *Database) GetTitles(ctx context.Context, offset, limit int) []Title {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "GetTitles")()

	res := []Title{}

	for i := db.lastTitleID - offset; i > db.lastTitleID-offset-limit; i-- {

		if title, ok := db.data.Titles[i]; ok {
			res = append(res, title.Copy(ctx))
		}
	}

	return res
}

func (db *Database) GetUnsuccessedPages(ctx context.Context) []PageFullInfo {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "GetUnsuccessPages")()

	res := []PageFullInfo{}

	for _, t := range db.data.Titles {
		for i, p := range t.Pages {
			if !p.Success {
				res = append(res, PageFullInfo{
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

func (db *Database) UpdateTitlePages(ctx context.Context, id int, pages []Page) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "UpdateTitlePages")()

	title, ok := db.data.Titles[id]
	if !ok {
		return TitleIndexError
	}

	title.Pages = pages
	title.Data.Parsed.Page = len(pages) > 0

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) GetPage(ctx context.Context, id, page int) (*PageFullInfo, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "GetPage")()

	title, ok := db.data.Titles[id]
	if !ok {
		return nil, TitleIndexError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return nil, PageIndexError
	}

	p := title.Pages[page]

	return &PageFullInfo{
		TitleID:    title.ID,
		PageNumber: page + 1,
		URL:        p.URL,
		Ext:        p.Ext,
		Success:    p.Success,
		LoadedAt:   p.LoadedAt,
	}, nil
}
