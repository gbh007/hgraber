package postgresql

import (
	"app/internal/domain/hgraber"
	"context"
	"database/sql"
)

func (d *Database) GetUnloadedBooks(ctx context.Context) []hgraber.Book {
	out := make([]hgraber.Book, 0)

	ids, err := d.bookUnprocessedMap(ctx)
	if err != nil {
		d.logger.Error(ctx, err)

		return out
	}

	for id := range ids {
		book, err := d.GetBook(ctx, id)
		if err != nil {
			d.logger.Error(ctx, err)
		} else {
			out = append(out, book)
		}
	}

	return out
}

func (d *Database) bookUnprocessedMap(ctx context.Context) (map[int]struct{}, error) {
	out := make(map[int]struct{})

	ids, err := d.bookUnprocessed(ctx)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		out[id] = struct{}{}
	}

	ids, err = d.bookAttrUnprocessed(ctx)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		out[id] = struct{}{}
	}

	return out, nil
}

func (d *Database) bookUnprocessed(ctx context.Context) ([]int, error) {
	ids := make([]int, 0)

	err := d.db.SelectContext(ctx, &ids, `SELECT id FROM books WHERE name IS NULL OR page_count IS NULL;`)
	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (d *Database) bookAttrUnprocessed(ctx context.Context) ([]int, error) {
	raw := make([]sql.NullInt64, 0)

	// TODO: проверить работоспособность.
	err := d.db.SelectContext(
		ctx, &raw,
		`SELECT l.book_id
		 FROM attributes AS a 
		 LEFT JOIN book_attributes_parsed AS l ON l.attr = a.code
		 WHERE l.parsed IS NULL OR l.parsed = FALSE
		 GROUP BY l.book_id;`)
	if err != nil {
		return nil, err
	}

	ids := make([]int, 0, len(raw))
	for _, id := range raw {
		ids = append(ids, int(id.Int64))
	}

	return ids, nil
}
