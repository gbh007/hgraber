package modelV2

import (
	"app/internal/domain/hgraber"
)

type RawBook struct {
	ID         int                     `json:"id"`
	Info       RawBookInfo             `json:"info"`
	Pages      []RawPage               `json:"pages,omitempty"`
	Attributes map[string]RawAttribute `json:"attributes,omitempty"`
}

func (book RawBook) Copy() RawBook {
	copyBook := RawBook{
		ID:         book.ID,
		Info:       book.Info.Copy(),
		Pages:      make([]RawPage, len(book.Pages)),
		Attributes: make(map[string]RawAttribute, len(book.Attributes)),
	}

	for i, page := range book.Pages {
		copyBook.Pages[i] = page.Copy()
	}

	for name, attr := range book.Attributes {
		copyBook.Attributes[name] = attr.Copy()
	}

	return copyBook
}

func (book RawBook) Super() hgraber.Book {
	domainBook := hgraber.Book{
		ID:      book.ID,
		Created: book.Info.Created,
		URL:     book.Info.URL,
		Pages:   make([]hgraber.Page, len(book.Pages)),
		Data: hgraber.BookInfo{
			Parsed: hgraber.BookInfoParsed{
				Name:       book.Info.Name != "",
				Page:       book.Info.PageCount > 0,
				Attributes: make(map[hgraber.Attribute]bool, len(book.Attributes)),
			},
			Name:       book.Info.Name,
			Rating:     book.Info.Rating,
			Attributes: make(map[hgraber.Attribute][]string, len(book.Attributes)),
		},
	}

	for code, attr := range book.Attributes {
		domainCode := hgraber.Attribute(code) // FIXME: по хорошему надо их матчить более явно

		domainBook.Data.Parsed.Attributes[domainCode] = attr.Parsed

		domainBook.Data.Attributes[domainCode] = make([]string, len(attr.Values))
		copy(domainBook.Data.Attributes[domainCode], attr.Values)
	}

	for i, p := range book.Pages {
		domainBook.Pages[i] = p.Super(book.ID)
	}

	return domainBook
}

func (tip RawBook) IsFullParsed() bool {
	if tip.Info.Name == "" || tip.Info.PageCount == 0 {
		return false
	}

	for _, attr := range tip.Attributes {
		if !attr.Parsed {
			return false
		}
	}

	return true
}

func RawBookFromDomain(book hgraber.Book) RawBook {
	rawBook := RawBook{
		ID: book.ID,
		Info: RawBookInfo{
			Created:   book.Created,
			URL:       book.URL,
			Name:      book.Data.Name,
			Rating:    book.Data.Rating,
			PageCount: len(book.Pages),
		},
		Pages:      RawPagesFromSuper(book.Pages),
		Attributes: make(map[string]RawAttribute, len(book.Data.Attributes)),
	}

	for code, attr := range book.Data.Attributes {
		values := make([]string, len(attr))
		copy(values, attr)

		rawBook.Attributes[string(code)] = RawAttribute{
			Parsed: book.Data.Parsed.Attributes[code],
			Values: values,
		}
	}

	return rawBook
}
