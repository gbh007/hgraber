package parser

import (
	"app/system"
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

func Load(ctx system.Context, URL string) (Parser, bool, error) {
	var (
		p   Parser
		err error
		ok  bool
	)
	switch {
	case strings.Index(URL, "https://imhentai.xxx/") == 0:
		p = &Parser_IMHENTAI_XXX{}
	case strings.Index(URL, "https://manga-online.biz/") == 0:
		p = &Parser_MANGAONLINE_BIZ{}
	default:
		err = fmt.Errorf("не корректная ссылка")
	}
	if err == nil {
		ok = p.Load(ctx, URL)
	}
	return p, ok, err
}

// Parser интерфейс для реализации парсеров для различных сайтов
type Parser interface {
	Load(ctx system.Context, URL string) bool
	ParseName(ctx system.Context) string
	ParsePages(ctx system.Context) []Page
	ParseTags(ctx system.Context) []string
	ParseAuthors(ctx system.Context) []string
	ParseCharacters(ctx system.Context) []string
	ParseLanguages(ctx system.Context) []string
	ParseCategories(ctx system.Context) []string
	ParseParodies(ctx system.Context) []string
	ParseGroups(ctx system.Context) []string
}

/*
func (p *Parser) Load(URL string) bool     { return false }
func (p Parser) ParseName(ctx system.Context) string         { return "" }
func (p Parser) ParsePages(ctx system.Context) []Page        { return []Page{} }
func (p Parser) ParseTags(ctx system.Context) []string       { return []string{} }
func (p Parser) ParseAuthors(ctx system.Context) []string    { return []string{} }
func (p Parser) ParseCharacters(ctx system.Context) []string { return []string{} }
func (p Parser) ParseLanguages(ctx system.Context) []string  { return []string{} }
func (p Parser) ParseCategories(ctx system.Context) []string { return []string{} }
func (p Parser) ParseParodies(ctx system.Context) []string   { return []string{} }
func (p Parser) ParseGroups(ctx system.Context) []string     { return []string{} }
*/
