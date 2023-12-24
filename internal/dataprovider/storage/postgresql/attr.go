package postgresql

import (
	"app/internal/domain/hgraber"
	"context"
)

var (
	attrMap = map[hgraber.Attribute]string{
		hgraber.AttrAuthor:    "author",
		hgraber.AttrCategory:  "category",
		hgraber.AttrCharacter: "character",
		hgraber.AttrGroup:     "group",
		hgraber.AttrLanguage:  "language",
		hgraber.AttrParody:    "parody",
		hgraber.AttrTag:       "tag",
	}

	reverseAttrMap = reverseAttr(attrMap)
)

func reverseAttr(in map[hgraber.Attribute]string) map[string]hgraber.Attribute {
	out := make(map[string]hgraber.Attribute, len(in))

	for k, v := range in {
		out[v] = k
	}

	return out
}

func (d *Database) UpdateAttributes(ctx context.Context, id int, attr hgraber.Attribute, data []string) error {
	attrCode, found := attrMap[attr]
	if !found {
		return hgraber.UnsupportedAttributeError
	}

	tx, err := d.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM book_attributes WHERE book_id = $1 AND attr = $2;`, id, attrCode)
	if err != nil {
		d.logger.IfErrFunc(ctx, tx.Rollback)

		return err
	}

	_, err = tx.ExecContext(ctx, `DELETE FROM book_attributes_parsed WHERE book_id = $1 AND attr = $2;`, id, attrCode)
	if err != nil {
		d.logger.IfErrFunc(ctx, tx.Rollback)

		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO book_attributes_parsed (book_id, attr, parsed) VALUES($1, $2, $3);`, id, attrCode, true)
	if err != nil {
		d.logger.IfErrFunc(ctx, tx.Rollback)

		return err
	}

	for _, v := range data {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO book_attributes (book_id, attr, value) VALUES($1, $2, $3);`,
			id, attrCode, v,
		)
		if err != nil {
			d.logger.IfErrFunc(ctx, tx.Rollback)

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

	err := d.db.SelectContext(ctx, &raw, `SELECT * FROM book_attributes WHERE book_id = $1;`, bookID)
	if err != nil {
		return nil, err
	}

	return raw, nil
}

func (d *Database) getBookAttrParsed(ctx context.Context, bookID int) ([]*BookAttributeParsed, error) {
	raw := make([]*BookAttributeParsed, 0)

	err := d.db.SelectContext(ctx, &raw, `SELECT * FROM book_attributes_parsed WHERE book_id = $1;`, bookID)
	if err != nil {
		return nil, err
	}

	return raw, nil
}
