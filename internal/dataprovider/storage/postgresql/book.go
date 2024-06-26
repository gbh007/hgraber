package postgresql

import (
	"app/internal/domain/hgraber"
	"context"
	"database/sql"
	"strings"
	"time"
)

func (d *Database) NewBook(ctx context.Context, name string, URL string, loaded bool) (int, error) {
	var count int
	err := d.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM books WHERE url = $1;`, URL)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, hgraber.BookAlreadyExistsError
	}

	var id int

	err = d.db.GetContext(
		ctx, &id,
		`INSERT INTO books (name, url, create_at) VALUES($1, $2, $3) RETURNING id;`,
		sql.NullString{String: name, Valid: loaded},
		strings.TrimSpace(URL),
		time.Now().UTC(),
	)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *Database) UpdateBookName(ctx context.Context, id int, name string) error {
	res, err := d.db.ExecContext(ctx, `UPDATE books SET name = $1 WHERE id = $2;`, name, id)
	if err != nil {
		return err
	}

	if !d.isApply(ctx, res) {
		return hgraber.BookNotFoundError
	}

	return nil
}

func (d *Database) UpdateBookRate(ctx context.Context, id int, rating int) error {
	res, err := d.db.ExecContext(ctx, `UPDATE books SET rate = $1 WHERE id = $2;`, rating, id)
	if err != nil {
		return err
	}

	if !d.isApply(ctx, res) {
		return hgraber.BookNotFoundError
	}

	return nil
}
