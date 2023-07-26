package model

import (
	"app/storage/schema"
	"context"
	"time"
)

type RawTitle struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
	URL     string    `json:"url"`

	Pages []RawPage    `json:"pages"`
	Data  RawTitleInfo `json:"info"`
}

func (t RawTitle) Copy(ctx context.Context) RawTitle {
	c := RawTitle{
		ID:      t.ID,
		Created: t.Created,
		URL:     t.URL,
		Pages:   make([]RawPage, len(t.Pages)),
		Data:    t.Data.Copy(ctx),
	}

	for i, p := range t.Pages {
		c.Pages[i] = p.Copy(ctx)
	}

	return c
}

func (t RawTitle) Super(ctx context.Context) schema.Title {
	c := schema.Title{
		ID:      t.ID,
		Created: t.Created,
		URL:     t.URL,
		Pages:   make([]schema.Page, len(t.Pages)),
		Data:    t.Data.Super(ctx),
	}

	for i, p := range t.Pages {
		c.Pages[i] = p.Super(ctx)
	}

	return c
}
