package model

import (
	"app/internal/domain"
	"context"
	"time"
)

type RawPage struct {
	URL      string    `json:"url"`
	Ext      string    `json:"ext"`
	Success  bool      `json:"success"`
	LoadedAt time.Time `json:"loaded_at"`
	Rate     int       `json:"rate,omitempty"`
}

func (p RawPage) Copy(ctx context.Context) RawPage {
	return RawPage{
		URL:      p.URL,
		Ext:      p.Ext,
		Success:  p.Success,
		LoadedAt: p.LoadedAt,
		Rate:     p.Rate,
	}
}

func (p RawPage) Super(ctx context.Context) domain.Page {
	return domain.Page{
		URL:      p.URL,
		Ext:      p.Ext,
		Success:  p.Success,
		LoadedAt: p.LoadedAt,
		Rate:     p.Rate,
	}
}

func RawPageFromSuper(p domain.Page) RawPage {
	return RawPage{
		URL:      p.URL,
		Ext:      p.Ext,
		Success:  p.Success,
		LoadedAt: p.LoadedAt,
		Rate:     p.Rate,
	}
}

func RawPagesFromSuper(in []domain.Page) []RawPage {
	out := make([]RawPage, len(in))

	for i, p := range in {
		out[i] = RawPageFromSuper(p)
	}

	return out
}
