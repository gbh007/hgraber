package jdb

import (
	"app/db"
	"app/system"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Page struct {
	URL      string    `json:"url"`
	Ext      string    `json:"ext"`
	Success  bool      `json:"success"`
	LoadedAt time.Time `json:"loaded_at"`
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

type TitlePages struct {
	Count int    `json:"count"`
	Pages []Page `json:"pages"`
}

type Title struct {
	ID      int       `json:"id"`
	Created time.Time `json:"created"`
	URL     string    `json:"url"`

	Pages TitlePages `json:"pages"`
	Data  TitleInfo  `json:"info"`
}

type DatabaseData struct {
	Titles map[int]Title `json:"titles"`
}

type Database struct {
	data  DatabaseData
	mutex *sync.RWMutex
}

func New() *Database {
	return &Database{
		mutex: &sync.RWMutex{},
		data: DatabaseData{
			Titles: make(map[int]Title),
		},
	}
}

func (dtb *Database) Save(ctx context.Context, path string) error {
	file, err := os.Create(path)
	defer system.IfErrFunc(ctx, file.Close)
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(dtb.data)
	if err != nil {
		system.Error(ctx, err)
		return err
	}

	return nil
}

func (dtb *Database) FetchFromSQL(ctx context.Context) {
	dtb.mutex.Lock()
	defer dtb.mutex.Unlock()

	checkPagesCount := db.SelectPagesCount(ctx)

	for _, title := range db.SelectTitles(ctx, 0, db.SelectTitlesCount(ctx)) {
		t := Title{
			ID:      title.ID,
			Created: title.Created,
			URL:     title.URL,
			Data: TitleInfo{
				Parsed: TitleInfoParsed{
					Name:       title.Loaded,
					Page:       title.ParsedPage,
					Tags:       title.ParsedTags,
					Authors:    title.ParsedAuthors,
					Characters: title.ParsedCharacters,
					Languages:  title.ParsedLanguages,
					Categories: title.ParsedCategories,
					Parodies:   title.ParsedParodies,
					Groups:     title.ParsedGroups,
				},
				Name:       title.Name,
				Tags:       title.Tags,
				Authors:    title.Authors,
				Characters: title.Characters,
				Languages:  title.Languages,
				Categories: title.Categories,
				Parodies:   title.Parodies,
				Groups:     title.Groups,
			},
			Pages: TitlePages{
				Count: title.PageCount,
				Pages: make([]Page, title.PageCount),
			},
		}

		for _, page := range db.SelectPagesByTitleID(ctx, t.ID) {
			t.Pages.Pages[page.PageNumber-1] = Page{
				URL:     page.URL,
				Ext:     page.Ext,
				Success: true,
			}
			checkPagesCount -= 1
		}

		dtb.data.Titles[t.ID] = t
	}

	if checkPagesCount != 0 {
		system.Error(ctx, fmt.Errorf("Ошибка конветриторования, несоответствие общего числа страниц"))
	}
}
