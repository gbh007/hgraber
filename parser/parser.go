package parser

import (
	"log"
	"strings"
)

func trimLastSlash(URL string, count int) string {
	c := 0
	ind := strings.LastIndexFunc(URL, func(r rune) bool {
		if r != rune('/') {
			return false
		}
		c++
		if c == count {
			return true
		}
		return false
	})
	return URL[:ind]
}

type Page struct {
	URL    string
	Number int
	Ext    string
}

func Load(URL string) (Parser, bool) {
	var p Parser
	switch {
	case strings.Index(URL, "https://imhentai.xxx/") == 0:
		p = &Parser_IMHENTAI_XXX{}
	case strings.Index(URL, "https://manga-online.biz/") == 0:
		p = &Parser_MANGAONLINE_BIZ{}
	default:
		log.Panicln("NO PARSER")
	}
	ok := p.Load(URL)
	return p, ok
}

// Parser интерфейс для реализации парсеров для различных сайтов
type Parser interface {
	Load(URL string) bool
	ParseName() string
	ParsePages() []Page
	ParseTags() []string
	ParseAuthors() []string
	ParseCharacters() []string
}

/*


func (p *Parser) Load(URL string) bool     { return false}
func (p Parser) ParseName() string         { return "" }
func (p Parser) ParsePages() []Page        { return []Page{} }
func (p Parser) ParseTags() []string       { return []string{} }
func (p Parser) ParseAuthors() []string    { return []string{} }
func (p Parser) ParseCharacters() []string { return []string{} }

*/
