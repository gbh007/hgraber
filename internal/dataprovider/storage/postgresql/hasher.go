package postgresql

import (
	"app/internal/domain/hgraber"
	"context"
	"database/sql"
)

func (d *Database) GetUnHashedPages(ctx context.Context) []hgraber.Page {
	raw := make([]*Page, 0)

	err := d.db.SelectContext(ctx, &raw, `SELECT * FROM pages WHERE "hash" IS NULL OR "size" IS NULL;`)
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

func (d *Database) UpdatePageHash(ctx context.Context, id int, page int, hash string, size int64) error {
	res, err := d.db.ExecContext(
		ctx,
		`UPDATE pages SET "hash" = $3, "size" = $4 WHERE book_id = $1 AND page_number = $2;`,
		id, page,
		sql.NullString{String: hash, Valid: hash != ""}, sql.NullInt64{Int64: size, Valid: size > 0},
	)
	if err != nil {
		return err
	}

	if !d.isApply(ctx, res) {
		return hgraber.PageNotFoundError
	}

	return nil
}
