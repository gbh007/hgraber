package rendering

import (
	"app/internal/domain"
	"fmt"
	"time"
)

func PageFromStorageWrap(addr string) func(raw domain.Page) Page {
	return func(raw domain.Page) Page {
		return Page{
			TitleID:    raw.BookID,
			PageNumber: raw.PageNumber,
			URL:        raw.URL,
			URLtoView:  fmt.Sprintf("%s/file/%d/%d.%s", addr, raw.BookID, raw.PageNumber, raw.Ext),
			Ext:        raw.Ext,
			Success:    raw.Success,
			LoadedAt:   raw.LoadedAt,
			Rate:       raw.Rate,
		}
	}
}

type Page struct {
	TitleID    int       `json:"title_id"`
	PageNumber int       `json:"page_number"`
	URL        string    `json:"url"`
	URLtoView  string    `json:"url_to_view"`
	Ext        string    `json:"ext"`
	Success    bool      `json:"success"`
	LoadedAt   time.Time `json:"loaded_at"`
	Rate       int       `json:"rate,omitempty"`
}

type TitleInfoParsed struct {
	Name       bool `json:"name,omitempty"`
	Page       bool `json:"page,omitempty"`
	Tags       bool `json:"tags,omitempty"`
	Authors    bool `json:"authors,omitempty"`
	Characters bool `json:"characters,omitempty"`
	Languages  bool `json:"languages,omitempty"`
	Categories bool `json:"categories,omitempty"`
	Parodies   bool `json:"parodies,omitempty"`
	Groups     bool `json:"groups,omitempty"`
}

func TitleInfoParsedFromStorage(raw domain.BookInfoParsed) TitleInfoParsed {
	return TitleInfoParsed{
		Name: raw.Name,
		Page: raw.Page,

		Tags:       raw.Attributes[domain.AttrTag],
		Authors:    raw.Attributes[domain.AttrAuthor],
		Characters: raw.Attributes[domain.AttrCharacter],
		Languages:  raw.Attributes[domain.AttrLanguage],
		Categories: raw.Attributes[domain.AttrCategory],
		Parodies:   raw.Attributes[domain.AttrParody],
		Groups:     raw.Attributes[domain.AttrGroup],
	}
}

type TitleInfo struct {
	Parsed     TitleInfoParsed `json:"parsed,omitempty"`
	Name       string          `json:"name,omitempty"`
	Rate       int             `json:"rate,omitempty"`
	Tags       []string        `json:"tags,omitempty"`
	Authors    []string        `json:"authors,omitempty"`
	Characters []string        `json:"characters,omitempty"`
	Languages  []string        `json:"languages,omitempty"`
	Categories []string        `json:"categories,omitempty"`
	Parodies   []string        `json:"parodies,omitempty"`
	Groups     []string        `json:"groups,omitempty"`
}

func TitleInfoFromStorage(raw domain.BookInfo) TitleInfo {
	out := TitleInfo{
		Parsed:     TitleInfoParsedFromStorage(raw.Parsed),
		Name:       raw.Name,
		Rate:       raw.Rate,
		Tags:       make([]string, len(raw.Attributes[domain.AttrTag])),
		Authors:    make([]string, len(raw.Attributes[domain.AttrAuthor])),
		Characters: make([]string, len(raw.Attributes[domain.AttrCharacter])),
		Languages:  make([]string, len(raw.Attributes[domain.AttrLanguage])),
		Categories: make([]string, len(raw.Attributes[domain.AttrCategory])),
		Parodies:   make([]string, len(raw.Attributes[domain.AttrParody])),
		Groups:     make([]string, len(raw.Attributes[domain.AttrGroup])),
	}

	copy(out.Tags, raw.Attributes[domain.AttrTag])
	copy(out.Authors, raw.Attributes[domain.AttrAuthor])
	copy(out.Characters, raw.Attributes[domain.AttrCharacter])
	copy(out.Languages, raw.Attributes[domain.AttrLanguage])
	copy(out.Categories, raw.Attributes[domain.AttrCategory])
	copy(out.Parodies, raw.Attributes[domain.AttrParody])
	copy(out.Groups, raw.Attributes[domain.AttrGroup])

	return out
}

type Title struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
	URL     string    `json:"url"`

	Pages []Page    `json:"pages"`
	Data  TitleInfo `json:"info"`
}

func TitleFromStorageWrap(addr string) func(raw domain.Book) Title {
	return func(raw domain.Book) Title {
		out := Title{
			ID:      raw.ID,
			Created: raw.Created,
			URL:     raw.URL,
			Pages:   make([]Page, len(raw.Pages)),
			Data:    TitleInfoFromStorage(raw.Data),
		}

		convertSlice(out.Pages, raw.Pages, PageFromStorageWrap(addr))

		return out
	}
}

func TitlesFromStorage(addr string, raw []domain.Book) []Title {
	out := make([]Title, len(raw))

	convertSlice(out, raw, TitleFromStorageWrap(addr))

	return out
}

type FirstHandleMultipleResult struct {
	TotalCount     int64    `json:"total_count"`
	LoadedCount    int64    `json:"loaded_count"`
	DuplicateCount int64    `json:"duplicate_count"`
	ErrorCount     int64    `json:"error_count"`
	NotHandled     []string `json:"not_handled,omitempty"`
}

func HandleMultipleResultFromDomain(raw domain.FirstHandleMultipleResult) FirstHandleMultipleResult {
	out := FirstHandleMultipleResult{
		TotalCount:     raw.TotalCount,
		LoadedCount:    raw.LoadedCount,
		DuplicateCount: raw.DuplicateCount,
		ErrorCount:     raw.ErrorCount,
		NotHandled:     make([]string, len(raw.NotHandled)),
	}

	copy(out.NotHandled, raw.NotHandled)

	return out
}
