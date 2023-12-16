package sqlite

import (
	"app/internal/domain"
	"app/system"
	"context"
	"database/sql"
	"errors"
)

func (d *Database) GetBooks(ctx context.Context, filter domain.BookFilter) []domain.Book {
	out := make([]domain.Book, 0)

	ids, err := d.bookIDs(ctx, filter)
	if err != nil {
		system.Error(ctx, err)

		return out
	}

	for _, id := range ids {
		book, err := d.GetBook(ctx, id)
		if err != nil {
			system.Error(ctx, err)
		} else {
			out = append(out, book)
		}
	}

	return out
}

func (d *Database) GetBook(ctx context.Context, bookID int) (domain.Book, error) {
	raw := new(Book)

	err := d.db.GetContext(ctx, raw, `SELECT * FROM books WHERE id = ? LIMIT 1;`, bookID)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Book{}, domain.BookNotFoundError
	}

	if err != nil {
		return domain.Book{}, err
	}

	out := domain.Book{
		ID:      raw.ID,
		Created: stringToTime(ctx, raw.CreateAt),
		URL:     raw.Url.String,
		Pages:   []domain.Page{},
		Data: domain.BookInfo{
			Parsed: domain.BookInfoParsed{
				Name:       raw.Name.Valid,
				Page:       raw.PageCount.Valid,
				Attributes: map[domain.Attribute]bool{},
			},
			Name:       raw.Name.String,
			Rate:       raw.Rate,
			Attributes: map[domain.Attribute][]string{},
		},
	}

	attributes, err := d.getBookAttr(ctx, bookID)
	if err != nil {
		return domain.Book{}, err
	}

	for _, attribute := range attributes {
		attr, ok := reverseAttrMap[attribute.Attr]
		if !ok {
			return domain.Book{}, domain.UnsupportedAttributeError
		}

		out.Data.Attributes[attr] = append(out.Data.Attributes[attr], attribute.Value)
	}

	attributesParsed, err := d.getBookAttrParsed(ctx, bookID)
	if err != nil {
		return domain.Book{}, err
	}

	for _, attribute := range attributesParsed {
		attr, ok := reverseAttrMap[attribute.Attr]
		if !ok {
			return domain.Book{}, domain.UnsupportedAttributeError
		}

		out.Data.Parsed.Attributes[attr] = attribute.Parsed
	}

	pages, err := d.getBookPages(ctx, bookID)
	if err != nil {
		return domain.Book{}, err
	}

	for _, p := range pages {
		out.Pages = append(out.Pages, pageToDomain(ctx, p))
	}

	return out, nil
}

func (d *Database) bookIDs(ctx context.Context, filter domain.BookFilter) ([]int, error) {
	ids := make([]int, 0)

	query := `SELECT id FROM books ORDER BY id ASC LIMIT ? OFFSET ?;`
	if filter.NewFirst {
		query = `SELECT id FROM books ORDER BY id DESC LIMIT ? OFFSET ?;`
	}

	err := d.db.SelectContext(ctx, &ids, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, err
	}

	return ids, nil
}
