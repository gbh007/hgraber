package model

import (
	"app/internal/domain"
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

func (t RawTitle) Super() domain.Title {
	c := domain.Title{
		ID:      t.ID,
		Created: t.Created,
		URL:     t.URL,
		Pages:   make([]domain.Page, len(t.Pages)),
		Data:    t.Data.Super(),
	}

	for i, p := range t.Pages {
		c.Pages[i] = p.Super()
	}

	return c
}
