package postgresql

import (
	"context"
)

func (d *Database) PagesCount(ctx context.Context) (c int) {
	err := d.db.GetContext(ctx, &c, `SELECT COUNT(*) FROM pages;`)
	if err != nil {
		d.logger.ErrorContext(ctx, err.Error())
	}

	return
}

func (d *Database) BooksCount(ctx context.Context) (c int) {
	err := d.db.GetContext(ctx, &c, `SELECT COUNT(*) FROM books;`)
	if err != nil {
		d.logger.ErrorContext(ctx, err.Error())
	}

	return
}

func (d *Database) UnloadedPagesCount(ctx context.Context) (c int) {
	err := d.db.GetContext(ctx, &c, `SELECT COUNT(*) FROM pages WHERE success = FALSE;`)
	if err != nil {
		d.logger.ErrorContext(ctx, err.Error())
	}

	return
}

func (d *Database) UnloadedBooksCount(ctx context.Context) int {
	m, err := d.bookUnprocessedMap(ctx)
	if err != nil {
		d.logger.ErrorContext(ctx, err.Error())
	}

	return len(m)
}
