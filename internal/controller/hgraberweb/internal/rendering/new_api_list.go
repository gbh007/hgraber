package rendering

import (
	"app/internal/domain/hgraber"
	"fmt"
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

	Rate        int      `json:"rate"`
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

func BookShortInfoFromDomain(addr string, book hgraber.Book) BookShortInfo {
	previewURL := ""
	pageCount := 0
	pageLoadedPercent := 0.0

	if len(book.Pages) > 0 {
		firstPage := book.Pages[0]
		pageCount = len(book.Pages)

		for _, page := range book.Pages {
			if page.Success {
				pageLoadedPercent++
			}
		}

		pageLoadedPercent = math.Round(pageLoadedPercent*10000/float64(pageCount)) / 100

		previewURL = fmt.Sprintf("%s/file/%d/%d.%s", addr, firstPage.BookID, firstPage.PageNumber, firstPage.Ext)
	}

	const renderTags = 8

	tags := make([]string, renderTags)
	tagCount := copy(tags, book.Data.Attributes[hgraber.AttrTag])
	tags = tags[:tagCount]
	hasMoreTags := len(book.Data.Attributes[hgraber.AttrTag]) > renderTags

	return BookShortInfo{
		ID:                book.ID,
		Created:           book.Created,
		PreviewURL:        previewURL,
		ParsedName:        book.Data.Parsed.Name,
		Name:              book.Data.Name,
		ParsedPage:        book.Data.Parsed.Page,
		PageCount:         pageCount,
		PageLoadedPercent: pageLoadedPercent,
		Rate:              book.Data.Rate,
		Tags:              tags,
		HasMoreTags:       hasMoreTags,
	}
}

func BookShortInfosFromDomain(addr string, books []hgraber.Book) []BookShortInfo {
	out := make([]BookShortInfo, len(books))

	convertSliceWithAddr(addr, out, books, BookShortInfoFromDomain)

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
