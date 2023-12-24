package jdb

import (
	"app/internal/domain/hgraber"
	"context"
	"time"
)

func (db *Database) UpdatePageSuccess(ctx context.Context, id, page int, success bool) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return hgraber.PageNotFoundError
	}

	title.Pages[page].Success = success
	if success {
		title.Pages[page].LoadedAt = time.Now()
	}

	db.data.Titles[id] = title
	db.needSave = true

	return nil
}

func (db *Database) UpdatePage(ctx context.Context, id int, page int, success bool, url string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	page--

	if page < 0 || page >= len(title.Pages) {
		return hgraber.PageNotFoundError
	}

	title.Pages[page].URL = url
	title.Pages[page].Success = success
	if success {
		title.Pages[page].LoadedAt = time.Now()
	}

	db.data.Titles[id] = title
	db.needSave = true

	return nil
}
