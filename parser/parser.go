package parser

import (
	"app/system/coreContext"
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

func Load(ctx coreContext.CoreContext, URL string) (Parser, bool, error) {
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
	Load(ctx coreContext.CoreContext, URL string) bool
	ParseName(ctx coreContext.CoreContext) string
	ParsePages(ctx coreContext.CoreContext) []Page
	ParseTags(ctx coreContext.CoreContext) []string
	ParseAuthors(ctx coreContext.CoreContext) []string
	ParseCharacters(ctx coreContext.CoreContext) []string
	ParseLanguages(ctx coreContext.CoreContext) []string
	ParseCategories(ctx coreContext.CoreContext) []string
	ParseParodies(ctx coreContext.CoreContext) []string
	ParseGroups(ctx coreContext.CoreContext) []string
}

/*
func (p *Parser) Load(URL string) bool     { return false }
func (p Parser) ParseName(ctx coreContext.CoreContext) string         { return "" }
func (p Parser) ParsePages(ctx coreContext.CoreContext) []Page        { return []Page{} }
func (p Parser) ParseTags(ctx coreContext.CoreContext) []string       { return []string{} }
func (p Parser) ParseAuthors(ctx coreContext.CoreContext) []string    { return []string{} }
func (p Parser) ParseCharacters(ctx coreContext.CoreContext) []string { return []string{} }
func (p Parser) ParseLanguages(ctx coreContext.CoreContext) []string  { return []string{} }
func (p Parser) ParseCategories(ctx coreContext.CoreContext) []string { return []string{} }
func (p Parser) ParseParodies(ctx coreContext.CoreContext) []string   { return []string{} }
func (p Parser) ParseGroups(ctx coreContext.CoreContext) []string     { return []string{} }
*/
