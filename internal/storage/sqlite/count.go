package sqlite

import (
	"app/system"
	"context"
)

func (d *Database) PagesCount(ctx context.Context) (c int) {
	err := d.db.GetContext(ctx, &c, `SELECT COUNT(*) FROM pages;`)
	system.IfErr(ctx, err)

	return
}

func (d *Database) TitlesCount(ctx context.Context) (c int) {
	err := d.db.GetContext(ctx, &c, `SELECT COUNT(*) FROM books;`)
	system.IfErr(ctx, err)

	return
}

func (d *Database) UnloadedPagesCount(ctx context.Context) (c int) {
	err := d.db.GetContext(ctx, &c, `SELECT COUNT(*) FROM pages WHERE success = FALSE;`)
	system.IfErr(ctx, err)

	return
}

func (d *Database) UnloadedTitlesCount(ctx context.Context) int {
	m, err := d.bookUnprocessedMap(ctx)
	system.IfErr(ctx, err)

	return len(m)
}
