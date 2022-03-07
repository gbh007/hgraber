package parser

import (
	"context"
	"fmt"
	"strings"
)

func trimLastSlash(URL string, count int) string {
	c := 0
	ind := strings.LastIndexFunc(URL, func(r rune) bool {
		if r != rune('/') {
			return false
		}
		c++
		return c == count
	})
	return URL[:ind]
}

type Page struct {
	URL    string
	Number int
	Ext    string
}

func Parse(ctx context.Context, URL string) (p Parser, err error) {
	switch {
	case strings.Index(URL, "https://imhentai.xxx/") == 0:
		p = &Parser_IMHENTAI_XXX{}
	case strings.Index(URL, "https://manga-online.biz/") == 0:
		p = &Parser_MANGAONLINE_BIZ{}
	default:
		err = fmt.Errorf("не корректная ссылка")
	}
	return p, err
}

func Load(ctx context.Context, URL string) (p Parser, ok bool, err error) {
	p, err = Parse(ctx, URL)
	if err == nil {
		ok = p.Load(ctx, URL)
	}
	return p, ok, err
}

// Parser интерфейс для реализации парсеров для различных сайтов
type Parser interface {
	Load(ctx context.Context, URL string) bool
	ParseName(ctx context.Context) string
	ParsePages(ctx context.Context) []Page
	ParseTags(ctx context.Context) []string
	ParseAuthors(ctx context.Context) []string
	ParseCharacters(ctx context.Context) []string
	ParseLanguages(ctx context.Context) []string
	ParseCategories(ctx context.Context) []string
	ParseParodies(ctx context.Context) []string
	ParseGroups(ctx context.Context) []string
}

/*
func (p *Parser) Load(URL string) bool     { return false }
func (p Parser) ParseName(ctx context.Context) string         { return "" }
func (p Parser) ParsePages(ctx context.Context) []Page        { return []Page{} }
func (p Parser) ParseTags(ctx context.Context) []string       { return []string{} }
func (p Parser) ParseAuthors(ctx context.Context) []string    { return []string{} }
func (p Parser) ParseCharacters(ctx context.Context) []string { return []string{} }
func (p Parser) ParseLanguages(ctx context.Context) []string  { return []string{} }
func (p Parser) ParseCategories(ctx context.Context) []string { return []string{} }
func (p Parser) ParseParodies(ctx context.Context) []string   { return []string{} }
func (p Parser) ParseGroups(ctx context.Context) []string     { return []string{} }
*/
