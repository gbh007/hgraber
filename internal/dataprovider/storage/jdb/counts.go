package jdb

import (
	"context"
)

func (db *Database) BooksCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	return len(db.data.Data.Books)
}

func (db *Database) UnloadedBooksCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	c := 0

	for _, t := range db.data.Data.Books {
		if !t.IsFullParsed() {
			c++
		}
	}

	return c
}

func (db *Database) PagesCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	c := 0

	for _, t := range db.data.Data.Books {
		c += len(t.Pages)
	}

	return c
}

func (db *Database) UnloadedPagesCount(ctx context.Context) int {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	c := 0

	for _, t := range db.data.Data.Books {
		for _, p := range t.Pages {
			if !p.Success {
				c++
			}
		}
	}

	return c
}

func (db *Database) PagesSize(ctx context.Context) (c int64) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	for _, t := range db.data.Data.Books {
		for _, p := range t.Pages {
			if p.Size > 0 {
				c += p.Size
			}
		}
	}

	return c
}
