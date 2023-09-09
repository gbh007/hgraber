package jdb

import (
	"context"
)

func (db *Database) TitlesCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	return len(db.data.Titles)
}

func (db *Database) UnloadedTitlesCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	c := 0

	for _, t := range db.data.Titles {
		if !t.Data.Parsed.IsFullParsed() {
			c++
		}
	}

	return c
}

func (db *Database) PagesCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	c := 0

	for _, t := range db.data.Titles {
		c += len(t.Pages)
	}

	return c
}

func (db *Database) UnloadedPagesCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	c := 0

	for _, t := range db.data.Titles {
		for _, p := range t.Pages {
			if !p.Success {
				c++
			}
		}
	}

	return c
}
