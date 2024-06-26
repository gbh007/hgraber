package postgresql

import (
	"app/internal/domain/hgraber"
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

func (d *Database) GetPage(ctx context.Context, id int, page int) (*hgraber.Page, error) {
	raw := new(Page)

	err := d.db.GetContext(
		ctx, raw,
		`SELECT * FROM pages WHERE book_id = $1 AND page_number = $2 LIMIT 1;`,
		id, page,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, hgraber.PageNotFoundError
	}

	if err != nil {
		return nil, err
	}

	p := pageToDomain(ctx, raw)

	return &p, nil
}

func (d *Database) GetUnsuccessPages(ctx context.Context) []hgraber.Page {
	raw := make([]*Page, 0)

	err := d.db.SelectContext(ctx, &raw, `SELECT * FROM pages WHERE success = FALSE;`)
	if err != nil {
		d.logger.ErrorContext(ctx, err.Error())

		return []hgraber.Page{}
	}

	out := make([]hgraber.Page, len(raw))
	for i, v := range raw {
		out[i] = pageToDomain(ctx, v)
	}

	return out
}

func (d *Database) UpdatePageSuccess(ctx context.Context, id int, page int, success bool) error {
	res, err := d.db.ExecContext(
		ctx,
		`UPDATE pages SET success = $1, load_at = $2 WHERE book_id = $3 AND page_number = $4;`,
		success, time.Now().UTC(), id, page,
	)
	if err != nil {
		return err
	}

	if !d.isApply(ctx, res) {
		return hgraber.PageNotFoundError
	}

	return nil
}

func (d *Database) UpdatePage(ctx context.Context, id int, page int, success bool, url string) error {
	res, err := d.db.ExecContext(
		ctx,
		`UPDATE pages SET success = $1, load_at = $2, url = $5 WHERE book_id = $3 AND page_number = $4;`,
		success, time.Now().UTC(), id, page, url,
	)
	if err != nil {
		return err
	}

	if !d.isApply(ctx, res) {
		return hgraber.PageNotFoundError
	}

	return nil
}

func (d *Database) UpdatePageRate(ctx context.Context, id int, page int, rating int) error {
	res, err := d.db.ExecContext(
		ctx,
		`UPDATE pages SET rate = $1 WHERE book_id = $2 AND page_number = $3;`,
		rating, id, page,
	)
	if err != nil {
		return err
	}

	if !d.isApply(ctx, res) {
		return hgraber.PageNotFoundError
	}

	return nil
}

func (d *Database) UpdateBookPages(ctx context.Context, id int, pages []hgraber.Page) error {
	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	res, err := tx.ExecContext(ctx, `UPDATE books SET page_count = $1 WHERE id = $2;`, len(pages), id)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			d.logger.ErrorContext(ctx, rollbackErr.Error())
		}

		return err
	}

	if !d.isApply(ctx, res) {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			d.logger.ErrorContext(ctx, rollbackErr.Error())
		}

		return hgraber.BookNotFoundError
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM pages WHERE book_id = $1;`, id)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			d.logger.ErrorContext(ctx, rollbackErr.Error())
		}

		return err
	}

	for _, v := range pages {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO pages (book_id, page_number, ext, url, success, create_at, load_at, rate, "hash", "size") VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
			id, v.PageNumber, v.Ext, strings.TrimSpace(v.URL), v.Success, time.Now().UTC(), sql.NullTime{Time: v.LoadedAt.UTC(), Valid: !v.LoadedAt.IsZero()}, v.Rating,
			sql.NullString{String: v.Hash, Valid: v.Hash != ""}, sql.NullInt64{Int64: v.Size, Valid: v.Size > 0},
		)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				d.logger.ErrorContext(ctx, rollbackErr.Error())
			}

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

	err := d.db.SelectContext(ctx, &raw, `SELECT * FROM pages WHERE book_id = $1 ORDER BY page_number;`, bookID)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func pageToDomain(_ context.Context, in *Page) hgraber.Page {
	return hgraber.Page{
		BookID:     in.BookID,
		PageNumber: in.PageNumber,
		URL:        in.Url,
		Ext:        in.Ext,
		Success:    in.Success,
		LoadedAt:   in.LoadAt.Time,
		Rating:     in.Rate,

		Hash: in.Hash.String,
		Size: in.Size.Int64,
	}
}
