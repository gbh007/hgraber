package model

import (
	"app/internal/domain/hgraber"
	"time"
)

type RawPage struct {
	URL      string    `json:"url"`
	Ext      string    `json:"ext"`
	Success  bool      `json:"success"`
	LoadedAt time.Time `json:"loaded_at"`
	Rate     int       `json:"rate,omitempty"`
}

func (p RawPage) Copy() RawPage {
	return RawPage{
		URL:      p.URL,
		Ext:      p.Ext,
		Success:  p.Success,
		LoadedAt: p.LoadedAt,
		Rate:     p.Rate,
	}
}

func (p RawPage) Super(bookID, pageNumber int) hgraber.Page {
	return hgraber.Page{
		BookID:     bookID,
		PageNumber: pageNumber,
		URL:        p.URL,
		Ext:        p.Ext,
		Success:    p.Success,
		LoadedAt:   p.LoadedAt,
		Rate:       p.Rate,
	}
}

func RawPageFromSuper(p hgraber.Page) RawPage {
	return RawPage{
		URL:      p.URL,
		Ext:      p.Ext,
		Success:  p.Success,
		LoadedAt: p.LoadedAt,
		Rate:     p.Rate,
	}
}

func RawPagesFromSuper(in []hgraber.Page) []RawPage {
	out := make([]RawPage, len(in))

	for i, p := range in {
		out[i] = RawPageFromSuper(p)
	}

	return out
}
