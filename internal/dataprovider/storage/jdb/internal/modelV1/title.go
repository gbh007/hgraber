package modelV1

import (
	"app/internal/domain/hgraber"
	"time"
)

type RawTitle struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
	URL     string    `json:"url"`

	Pages []RawPage    `json:"pages"`
	Data  RawTitleInfo `json:"info"`
}

func (t RawTitle) Copy() RawTitle {
	c := RawTitle{
		ID:      t.ID,
		Created: t.Created,
		URL:     t.URL,
		Pages:   make([]RawPage, len(t.Pages)),
		Data:    t.Data.Copy(),
	}

	for i, p := range t.Pages {
		c.Pages[i] = p.Copy()
	}

	return c
}

func (t RawTitle) Super() hgraber.Book {
	c := hgraber.Book{
		ID:      t.ID,
		Created: t.Created,
		URL:     t.URL,
		Pages:   make([]hgraber.Page, len(t.Pages)),
		Data:    t.Data.Super(),
	}

	for i, p := range t.Pages {
		c.Pages[i] = p.Super(t.ID, i+1)
	}

	return c
}
