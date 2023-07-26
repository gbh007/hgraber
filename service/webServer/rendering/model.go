package rendering

import (
	"app/storage/schema"
	"time"
)

type Page struct {
	URL      string    `json:"url"`
	Ext      string    `json:"ext"`
	Success  bool      `json:"success"`
	LoadedAt time.Time `json:"loaded_at"`
	Rate     int       `json:"rate,omitempty"`
}

func PageFromStorage(raw schema.Page) Page {
	return Page{
		URL:      raw.URL,
		Ext:      raw.Ext,
		Success:  raw.Success,
		LoadedAt: raw.LoadedAt,
		Rate:     raw.Rate,
	}
}

type PageFullInfo struct {
	TitleID    int       `json:"title_id"`
	PageNumber int       `json:"page_number"`
	URL        string    `json:"url"`
	Ext        string    `json:"ext"`
	Success    bool      `json:"success"`
	LoadedAt   time.Time `json:"loaded_at"`
	Rate       int       `json:"rate,omitempty"`
}

func PageFullInfoFromStorage(raw *schema.PageFullInfo) PageFullInfo {
	return PageFullInfo{
		TitleID:    raw.TitleID,
		PageNumber: raw.PageNumber,
		URL:        raw.URL,
		Ext:        raw.Ext,
		Success:    raw.Success,
		LoadedAt:   raw.LoadedAt,
		Rate:       raw.Rate,
	}
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

func TitleInfoParsedFromStorage(raw schema.TitleInfoParsed) TitleInfoParsed {
	return TitleInfoParsed{
		Name:       raw.Name,
		Page:       raw.Page,
		Tags:       raw.Tags,
		Authors:    raw.Authors,
		Characters: raw.Characters,
		Languages:  raw.Languages,
		Categories: raw.Categories,
		Parodies:   raw.Parodies,
		Groups:     raw.Groups,
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

func TitleInfoFromStorage(raw schema.TitleInfo) TitleInfo {
	out := TitleInfo{
		Parsed:     TitleInfoParsedFromStorage(raw.Parsed),
		Name:       raw.Name,
		Rate:       raw.Rate,
		Tags:       make([]string, len(raw.Tags)),
		Authors:    make([]string, len(raw.Authors)),
		Characters: make([]string, len(raw.Characters)),
		Languages:  make([]string, len(raw.Languages)),
		Categories: make([]string, len(raw.Categories)),
		Parodies:   make([]string, len(raw.Parodies)),
		Groups:     make([]string, len(raw.Groups)),
	}

	copy(out.Tags, raw.Tags)
	copy(out.Authors, raw.Authors)
	copy(out.Characters, raw.Characters)
	copy(out.Languages, raw.Languages)
	copy(out.Categories, raw.Categories)
	copy(out.Parodies, raw.Parodies)
	copy(out.Groups, raw.Groups)

	return out
}

type Title struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
	URL     string    `json:"url"`

	Pages []Page    `json:"pages"`
	Data  TitleInfo `json:"info"`
}

func TitleFromStorage(raw schema.Title) Title {
	out := Title{
		ID:      raw.ID,
		Created: raw.Created,
		URL:     raw.URL,
		Pages:   make([]Page, len(raw.Pages)),
		Data:    TitleInfoFromStorage(raw.Data),
	}

	for index, page := range raw.Pages {
		out.Pages[index] = PageFromStorage(page)
	}

	return out
}
func TitlesFromStorage(raw []schema.Title) []Title {
	out := make([]Title, 0, len(raw))

	for _, t := range raw {
		out = append(out, TitleFromStorage(t))
	}

	return out
}
