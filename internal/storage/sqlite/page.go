package sqlite

import (
	"app/internal/domain"
	"app/system"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

func (d *Database) GetPage(ctx context.Context, id int, page int) (*domain.PageFullInfo, error) {
	raw := new(Page)

	err := d.db.GetContext(
		ctx, raw,
		`SELECT * FROM pages WHERE book_id = ? AND page_number = ? LIMIT 1;`,
		id, page,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.PageNotFoundError
	}

	if err != nil {
		return nil, err
	}

	p := pageToDomain(ctx, raw)

	return &p, nil
}

func (d *Database) GetUnsuccessPages(ctx context.Context) []domain.PageFullInfo {
	raw := make([]*Page, 0)

	err := d.db.SelectContext(ctx, &raw, `SELECT * FROM pages WHERE success = FALSE;`)
	if err != nil {
		system.Error(ctx, err)

		return []domain.PageFullInfo{}
	}

	out := make([]domain.PageFullInfo, len(raw))
	for i, v := range raw {
		out[i] = pageToDomain(ctx, v)
	}

	return out
}

func (d *Database) UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error {
	res, err := d.db.ExecContext(
		ctx,
		`UPDATE pages SET success = ?, load_at = ? WHERE book_id = ? AND page_number = ?;`,
		success, timeToSQLString(time.Now()), id, page,
	)
	if err != nil {
		return err
	}

	if !isApply(ctx, res) {
		return domain.PageNotFoundError
	}

	return nil

}

func (d *Database) UpdatePageRate(ctx context.Context, id int, page int, rate int) error {
	res, err := d.db.ExecContext(
		ctx,
		`UPDATE pages SET rate = ? WHERE book_id = ? AND page_number = ?;`,
		rate, id, page,
	)
	if err != nil {
		return err
	}

	if !isApply(ctx, res) {
		return domain.PageNotFoundError
	}

	return nil
}

func (d *Database) UpdateBookPages(ctx context.Context, id int, pages []domain.Page) error {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE books SET page_count = ? WHERE id = ?;`, len(pages), id)
	if err != nil {
		system.IfErrFunc(ctx, tx.Rollback)

		return err
	}

	if !isApply(ctx, res) {
		system.IfErrFunc(ctx, tx.Rollback)

		return domain.BookNotFoundError
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM pages WHERE book_id = ?;`, id)
	if err != nil {
		system.IfErrFunc(ctx, tx.Rollback)

		return err
	}

	for i, v := range pages {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO pages (book_id, page_number, ext, url, success, create_at, load_at, rate) VALUES(?, ?, ?, ?, ?, ?, ?, ?);`,
			id, i+1, v.Ext, strings.TrimSpace(v.URL), v.Success, timeToString(time.Now()), timeToSQLString(v.LoadedAt), v.Rate,
		)
		if err != nil {
			system.IfErrFunc(ctx, tx.Rollback)

			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) getBookPages(ctx context.Context, bookID int) ([]*Page, error) {
	raw := make([]*Page, 0)

	err := d.db.SelectContext(ctx, &raw, `SELECT * FROM pages WHERE book_id = ? ORDER BY page_number;`, bookID)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func pageToDomain(ctx context.Context, in *Page) domain.PageFullInfo {
	return domain.PageFullInfo{
		BookID:     in.BookID,
		PageNumber: in.PageNumber,
		URL:        in.Url,
		Ext:        in.Ext,
		Success:    in.Success,
		LoadedAt:   stringToTime(ctx, in.LoadAt.String),
		Rate:       in.Rate,
	}
}
