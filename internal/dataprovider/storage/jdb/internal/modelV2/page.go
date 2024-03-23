package modelV2

import (
	"app/internal/domain/hgraber"
	"time"
)

type RawPage struct {
	PageNumber int       `json:"page_number"`
	URL        string    `json:"url"`
	Ext        string    `json:"ext"`
	Success    bool      `json:"success"`
	LoadAt     time.Time `json:"load_at"`
	Rating     int       `json:"rating,omitempty"`

	Hash string `json:"hash,omitempty"`
	Size int64  `json:"size,omitempty"`
}

func (p RawPage) Copy() RawPage {
	return RawPage{
		PageNumber: p.PageNumber,
		URL:        p.URL,
		Ext:        p.Ext,
		Success:    p.Success,
		LoadAt:     p.LoadAt,
		Rating:     p.Rating,
		Hash:       p.Hash,
		Size:       p.Size,
	}
}

func (p RawPage) Super(bookID int) hgraber.Page {
	return hgraber.Page{
		BookID:     bookID,
		PageNumber: p.PageNumber,
		URL:        p.URL,
		Ext:        p.Ext,
		Success:    p.Success,
		LoadedAt:   p.LoadAt,
		Rating:     p.Rating,
		Hash:       p.Hash,
		Size:       p.Size,
	}
}

func RawPageFromSuper(p hgraber.Page) RawPage {
	return RawPage{
		PageNumber: p.PageNumber,
		URL:        p.URL,
		Ext:        p.Ext,
		Success:    p.Success,
		LoadAt:     p.LoadedAt,
		Rating:     p.Rating,
		Hash:       p.Hash,
		Size:       p.Size,
	}
}

func RawPagesFromSuper(in []hgraber.Page) []RawPage {
	out := make([]RawPage, len(in))

	for i, p := range in {
		out[i] = RawPageFromSuper(p)
	}

	return out
}
