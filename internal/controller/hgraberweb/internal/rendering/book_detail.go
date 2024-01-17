package rendering

import (
	"app/internal/domain/hgraber"
	"math"
	"slices"
	"time"
)

type BookDetailInfo struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`

	PreviewURL string `json:"preview_url,omitempty"`

	ParsedName bool   `json:"parsed_name"`
	Name       string `json:"name"`

	ParsedPage        bool    `json:"parsed_page"`
	PageCount         int     `json:"page_count"`
	PageLoadedPercent float64 `json:"page_loaded_percent"`

	Rating int `json:"rating"`

	Attributes []BookDetailAttributeInfo `json:"attributes,omitempty"`
	Pages      []BookDetailPagePreview   `json:"pages,omitempty"`
}

type BookDetailAttributeInfo struct {
	Name   string   `json:"name"`
	Values []string `json:"values"`
	order  int      `json:"-"` // FIXME: не самое подходящая реализация
}

type BookDetailPagePreview struct {
	PageNumber int    `json:"page_number"`
	PreviewURL string `json:"preview_url,omitempty"`
	Rating     int    `json:"rating"`
}

func BookDetailPagePreviewFromDomain(addr string, raw hgraber.Page) BookDetailPagePreview {
	previewURL := ""

	if raw.Success {
		previewURL = fileURL(addr, raw.BookID, raw.PageNumber, raw.Ext)
	}

	return BookDetailPagePreview{
		PageNumber: raw.PageNumber,
		PreviewURL: previewURL,
		Rating:     raw.Rating,
	}
}

func BookDetailInfoFromDomain(addr string, raw hgraber.Book) BookDetailInfo {
	previewURL := ""
	pageCount := 0
	pageLoadedPercent := 0.0

	if len(raw.Pages) > 0 {
		firstPage := raw.Pages[0]
		pageCount = len(raw.Pages)

		for _, page := range raw.Pages {
			if page.Success {
				pageLoadedPercent++
			}
		}

		pageLoadedPercent = math.Round(pageLoadedPercent*10000/float64(pageCount)) / 100

		previewURL = fileURL(addr, firstPage.BookID, firstPage.PageNumber, firstPage.Ext)
	}

	attrs := make([]BookDetailAttributeInfo, 0, len(raw.Data.Attributes))
	for code, attr := range raw.Data.Attributes {
		if len(attr) < 1 {
			continue
		}

		values := make([]string, len(attr))
		copy(values, attr)

		attrs = append(attrs, BookDetailAttributeInfo{
			Name:   attributeDisplayName(code),
			Values: values,
			order:  attributeOrder(code),
		})
	}

	slices.SortFunc(attrs, func(a, b BookDetailAttributeInfo) int {
		return a.order - b.order
	})

	pages := make([]BookDetailPagePreview, len(raw.Pages))
	convertSliceWithAddr(addr, pages, raw.Pages, BookDetailPagePreviewFromDomain)

	return BookDetailInfo{
		ID:                raw.ID,
		Created:           raw.Created,
		PreviewURL:        previewURL,
		ParsedName:        raw.Data.Parsed.Name,
		Name:              raw.Data.Name,
		ParsedPage:        raw.Data.Parsed.Page,
		PageCount:         pageCount,
		PageLoadedPercent: pageLoadedPercent,
		Rating:            raw.Data.Rating,
		Attributes:        attrs,
		Pages:             pages,
	}
}
