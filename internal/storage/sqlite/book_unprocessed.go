package sqlite

import (
	"app/internal/domain"
	"app/system"
	"context"
)

func (d *Database) GetUnloadedTitles(ctx context.Context) []domain.Title {
	out := make([]domain.Title, 0)

	ids, err := d.bookUnprocessedMap(ctx)
	if err != nil {
		system.Error(ctx, err)

		return out
	}

	for id := range ids {
		book, err := d.GetTitle(ctx, id)
		if err != nil {
			system.Error(ctx, err)
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
	ids := make([]int, 0)

	// TODO: проверить работоспособность.
	err := d.db.SelectContext(
		ctx, &ids,
		`SELECT b.id 
		 FROM books as b 
		 LEFT JOIN attributes as a 
		 LEFT JOIN book_attributes_parsed as l ON l.book_id = b.id AND l.attr = a.code
		 WHERE l.parsed IS NULL OR l.parsed = FALSE
		 GROUP BY b.id;`)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
