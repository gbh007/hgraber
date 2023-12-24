package postgresql

import (
	"app/internal/domain/hgraber"
	"context"
	"database/sql"
	"errors"
)

func (d *Database) GetBooks(ctx context.Context, filter hgraber.BookFilter) []hgraber.Book {
	out := make([]hgraber.Book, 0)

	ids, err := d.bookIDs(ctx, filter)
	if err != nil {
		d.logger.Error(ctx, err)

		return out
	}

	for _, id := range ids {
		book, err := d.GetBook(ctx, id)
		if err != nil {
			d.logger.Error(ctx, err)
		} else {
			out = append(out, book)
		}
	}

	return out
}

func (d *Database) GetBook(ctx context.Context, bookID int) (hgraber.Book, error) {
	raw := new(Book)

	err := d.db.GetContext(ctx, raw, `SELECT * FROM books WHERE id = $1 LIMIT 1;`, bookID)
	if errors.Is(err, sql.ErrNoRows) {
		return hgraber.Book{}, hgraber.BookNotFoundError
	}

	if err != nil {
		return hgraber.Book{}, err
	}

	out := hgraber.Book{
		ID:      raw.ID,
		Created: raw.CreateAt,
		URL:     raw.Url.String,
		Pages:   []hgraber.Page{},
		Data: hgraber.BookInfo{
			Parsed: hgraber.BookInfoParsed{
				Name:       raw.Name.Valid,
				Page:       raw.PageCount.Valid,
				Attributes: map[hgraber.Attribute]bool{},
			},
			Name:       raw.Name.String,
			Rate:       raw.Rate,
			Attributes: map[hgraber.Attribute][]string{},
		},
	}

	attributes, err := d.getBookAttr(ctx, bookID)
	if err != nil {
		return hgraber.Book{}, err
	}

	for _, attribute := range attributes {
		attr, ok := reverseAttrMap[attribute.Attr]
		if !ok {
			return hgraber.Book{}, hgraber.UnsupportedAttributeError
		}

		out.Data.Attributes[attr] = append(out.Data.Attributes[attr], attribute.Value)
	}

	attributesParsed, err := d.getBookAttrParsed(ctx, bookID)
	if err != nil {
		return hgraber.Book{}, err
	}

	for _, attribute := range attributesParsed {
		attr, ok := reverseAttrMap[attribute.Attr]
		if !ok {
			return hgraber.Book{}, hgraber.UnsupportedAttributeError
		}

		out.Data.Parsed.Attributes[attr] = attribute.Parsed
	}

	pages, err := d.getBookPages(ctx, bookID)
	if err != nil {
		return hgraber.Book{}, err
	}

	for _, p := range pages {
		out.Pages = append(out.Pages, pageToDomain(ctx, p))
	}

	return out, nil
}

func (d *Database) bookIDs(ctx context.Context, filter hgraber.BookFilter) ([]int, error) {
	ids := make([]int, 0)

	query := `SELECT id FROM books ORDER BY id ASC LIMIT $1 OFFSET $2;`
	if filter.NewFirst {
		query = `SELECT id FROM books ORDER BY id DESC LIMIT $1 OFFSET $2;`
	}

	err := d.db.SelectContext(ctx, &ids, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
