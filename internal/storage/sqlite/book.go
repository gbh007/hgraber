package sqlite

import (
	"app/internal/domain"
	"context"
	"database/sql"
	"strings"
	"time"
)

func (d *Database) NewBook(ctx context.Context, name string, URL string, loaded bool) (int, error) {
	var count int
	err := d.db.GetContext(ctx, &count, `SELECT COUNT(*) FROM books WHERE url = ?;`, URL)
	if err != nil {
		return 0, err
	}

	if count > 0 {
		return 0, domain.BookAlreadyExistsError
	}

	res, err := d.db.ExecContext(
		ctx,
		`INSERT INTO books (name, url, create_at) VALUES(?, ?, ?);`,
		sql.NullString{String: name, Valid: loaded},
		strings.TrimSpace(URL),
		timeToString(time.Now()),
	)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (d *Database) UpdateBookName(ctx context.Context, id int, name string) error {
	res, err := d.db.ExecContext(ctx, `UPDATE books SET name = ? WHERE id = ?;`, name, id)
	if err != nil {
		return err
	}

	if !isApply(ctx, res) {
		return domain.BookNotFoundError
	}

	return nil
}

func (d *Database) UpdateBookRate(ctx context.Context, id int, rate int) error {
	res, err := d.db.ExecContext(ctx, `UPDATE books SET rate = ? WHERE id = ?;`, rate, id)
	if err != nil {
		return err
	}

	if !isApply(ctx, res) {
		return domain.BookNotFoundError
	}

	return nil
}
