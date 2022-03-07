package jdb

import (
	"app/system"
	"context"
	"time"
)

var copyStopwatch = false

func DisableCopyStopwatch(ctx context.Context) {
	copyStopwatch = false
	system.Warning(ctx, "Отключен режим отладки копирования")
}

func EnableCopyStopwatch(ctx context.Context) {
	copyStopwatch = true
	system.Warning(ctx, "Включен режим отладки копирования")
}

func CopyStopwatchStatus() bool {
	return copyStopwatch
}

type Page struct {
	URL      string    `json:"url"`
	Ext      string    `json:"ext"`
	Success  bool      `json:"success"`
	LoadedAt time.Time `json:"loaded_at"`
}

func (p Page) Copy(ctx context.Context) Page {
	if copyStopwatch {
		defer system.Stopwatch(ctx, "Скопизованы данные о странице за")()
	}

	return Page{
		URL:      p.URL,
		Ext:      p.Ext,
		Success:  p.Success,
		LoadedAt: p.LoadedAt,
	}
}

type PageFullInfo struct {
	TitleID    int       `json:"title_id"`
	PageNumber int       `json:"page_number"`
	URL        string    `json:"url"`
	Ext        string    `json:"ext"`
	Success    bool      `json:"success"`
	LoadedAt   time.Time `json:"loaded_at"`
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

func (tip TitleInfoParsed) Copy(ctx context.Context) TitleInfoParsed {
	if copyStopwatch {
		defer system.Stopwatch(ctx, "Скопизованы данные о парсе тайтла за")()
	}

	return TitleInfoParsed{
		Name:       tip.Name,
		Page:       tip.Page,
		Tags:       tip.Tags,
		Authors:    tip.Authors,
		Characters: tip.Characters,
		Languages:  tip.Languages,
		Categories: tip.Categories,
		Parodies:   tip.Parodies,
		Groups:     tip.Groups,
	}
}

func (tip TitleInfoParsed) IsFullParsed(ctx context.Context) bool {
	return tip.Name &&
		tip.Page &&
		tip.Tags &&
		tip.Authors &&
		tip.Characters &&
		tip.Languages &&
		tip.Categories &&
		tip.Parodies &&
		tip.Groups
}

type TitleInfo struct {
	Parsed     TitleInfoParsed `json:"parsed,omitempty"`
	Name       string          `json:"name,omitempty"`
	Tags       []string        `json:"tags,omitempty"`
	Authors    []string        `json:"authors,omitempty"`
	Characters []string        `json:"characters,omitempty"`
	Languages  []string        `json:"languages,omitempty"`
	Categories []string        `json:"categories,omitempty"`
	Parodies   []string        `json:"parodies,omitempty"`
	Groups     []string        `json:"groups,omitempty"`
}

func (ti TitleInfo) Copy(ctx context.Context) TitleInfo {
	if copyStopwatch {
		defer system.Stopwatch(ctx, "Скопированы данные информации о тайтле за")()
	}

	c := TitleInfo{
		Parsed:     ti.Parsed.Copy(ctx),
		Name:       ti.Name,
		Tags:       make([]string, len(ti.Tags)),
		Authors:    make([]string, len(ti.Authors)),
		Characters: make([]string, len(ti.Characters)),
		Languages:  make([]string, len(ti.Languages)),
		Categories: make([]string, len(ti.Categories)),
		Parodies:   make([]string, len(ti.Parodies)),
		Groups:     make([]string, len(ti.Groups)),
	}

	copy(c.Tags, ti.Tags)
	copy(c.Authors, ti.Authors)
	copy(c.Characters, ti.Characters)
	copy(c.Languages, ti.Languages)
	copy(c.Categories, ti.Categories)
	copy(c.Parodies, ti.Parodies)
	copy(c.Groups, ti.Groups)

	return c
}

type Title struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
	URL     string    `json:"url"`

	Pages []Page    `json:"pages"`
	Data  TitleInfo `json:"info"`
}

func (t Title) Copy(ctx context.Context) Title {
	if copyStopwatch {
		defer system.Stopwatch(ctx, "Скопированы данные о тайтле за")()
	}

	c := Title{
		ID:      t.ID,
		Created: t.Created,
		URL:     t.URL,
		Pages:   make([]Page, len(t.Pages)),
		Data:    t.Data.Copy(ctx),
	}

	for i, p := range t.Pages {
		c.Pages[i] = p.Copy(ctx)
	}

	return c
}
