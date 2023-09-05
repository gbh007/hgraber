package sqlite

import (
	"app/internal/domain"
	"app/system"
	"context"
)

var (
	attrMap = map[domain.Attribute]string{
		domain.AttrAuthor:    "author",
		domain.AttrCategory:  "category",
		domain.AttrCharacter: "character",
		domain.AttrGroup:     "group",
		domain.AttrLanguage:  "language",
		domain.AttrParody:    "parody",
		domain.AttrTag:       "tag",
	}

	reverseAttrMap = reverseAttr(attrMap)
)

func reverseAttr(in map[domain.Attribute]string) map[string]domain.Attribute {
	out := make(map[string]domain.Attribute, len(in))

	for k, v := range in {
		out[v] = k
	}

	return out
}

func (d *Database) UpdateAttributes(ctx context.Context, id int, attr domain.Attribute, data []string) error {
	attrCode, found := attrMap[attr]
	if !found {
		return domain.UnsupportedAttributeError
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM book_attributes WHERE book_id = ? AND attr = ?;`, id, attrCode)
	if err != nil {
		system.IfErrFunc(ctx, tx.Rollback)

		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM book_attributes_parsed WHERE book_id = ? AND attr = ?;`, id, attrCode)
	if err != nil {
		system.IfErrFunc(ctx, tx.Rollback)

		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO book_attributes_parsed (book_id, attr, parsed) VALUES(?, ?, ?);`, id, attrCode, true)
	if err != nil {
		system.IfErrFunc(ctx, tx.Rollback)

		return err
	}

	for _, v := range data {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO book_attributes (book_id, attr, value) VALUES(?, ?, ?);`,
			id, attrCode, v,
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

func (d *Database) getBookAttr(ctx context.Context, bookID int) ([]*BookAttribute, error) {
	raw := make([]*BookAttribute, 0)

	err := d.db.SelectContext(ctx, &raw, `SELECT * FROM book_attributes WHERE book_id = ?;`, bookID)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (d *Database) getBookAttrParsed(ctx context.Context, bookID int) ([]*BookAttributeParsed, error) {
	raw := make([]*BookAttributeParsed, 0)

	err := d.db.SelectContext(ctx, &raw, `SELECT * FROM book_attributes_parsed WHERE book_id = ?;`, bookID)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
