package rendering

import (
	"app/internal/domain/hgraber"
	"math"
	"time"
)

type BookShortInfo struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`

	PreviewURL string `json:"preview_url,omitempty"`

	ParsedName bool   `json:"parsed_name"`
	Name       string `json:"name"`

	ParsedPage        bool    `json:"parsed_page"`
	PageCount         int     `json:"page_count"`
	PageLoadedPercent float64 `json:"page_loaded_percent"`

	Rating      int      `json:"rating"`
	Tags        []string `json:"tags,omitempty"`
	HasMoreTags bool     `json:"has_more_tags"`
}

type PageForPagination struct {
	Value       int  `json:"value"`
	IsCurrent   bool `json:"is_current"`
	IsSeparator bool `json:"is_separator"`
}

type BookListResponse struct {
	Books []BookShortInfo     `json:"books"`
	Pages []PageForPagination `json:"pages"`
}

func BookShortInfoFromDomain(addr string, raw hgraber.Book) BookShortInfo {
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

	const renderTags = 8

	tags := make([]string, renderTags)
	tagCount := copy(tags, raw.Data.Attributes[hgraber.AttrTag])
	tags = tags[:tagCount]
	hasMoreTags := len(raw.Data.Attributes[hgraber.AttrTag]) > renderTags

	return BookShortInfo{
		ID:                raw.ID,
		Created:           raw.Created,
		PreviewURL:        previewURL,
		ParsedName:        raw.Data.Parsed.Name,
		Name:              raw.Data.Name,
		ParsedPage:        raw.Data.Parsed.Page,
		PageCount:         pageCount,
		PageLoadedPercent: pageLoadedPercent,
		Rating:            raw.Data.Rating,
		Tags:              tags,
		HasMoreTags:       hasMoreTags,
	}
}

func BookShortInfosFromDomain(addr string, raw []hgraber.Book) []BookShortInfo {
	out := make([]BookShortInfo, len(raw))

	convertSliceWithAddr(addr, out, raw, BookShortInfoFromDomain)

	return out
}

func BookListResponseFromDomain(addr string, raw hgraber.FilteredBooks) BookListResponse {
	pages := make([]PageForPagination, 0, len(raw.Pages))

	for _, page := range raw.Pages {
		pages = append(pages, PageForPagination{
			Value:       page,
			IsCurrent:   page == raw.CurrentPage,
			IsSeparator: page == -1, // TODO: заменить на менее костыльный вариант
		})
	}

	return BookListResponse{
		Books: BookShortInfosFromDomain(addr, raw.Books),
		Pages: pages,
	}
}
