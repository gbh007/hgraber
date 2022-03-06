package jdb

import (
	"app/system"
	"context"
)

/*

SelectTitlesCount
SelectUnloadTitlesCount
SelectPagesCount
SelectUnloadPagesCount


*/

func (db *Database) TitlesCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "TitlesCount")()

	return len(db.data.Titles)
}

func (db *Database) UnloadedTitlesCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "UnloadedTitlesCount")()

	c := 0

	for _, t := range db.data.Titles {
		if !t.Data.Parsed.IsFullParsed(ctx) {
			c++
		}
	}

	return c
}

func (db *Database) PagesCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "PagesCount")()

	c := 0

	for _, t := range db.data.Titles {
		c += len(t.Pages)
	}

	return c
}

func (db *Database) UnloadedPagesCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	defer system.Stopwatch(ctx, "UnloadedPagesCount")()

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
