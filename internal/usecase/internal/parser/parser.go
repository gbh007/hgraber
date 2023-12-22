package parser

import (
	"app/internal/domain"
	"context"
	"errors"
	"strings"
)

var ErrInvalidLink = errors.New("не корректная ссылка")

type Requester interface {
	RequestString(ctx context.Context, URL string) (string, error)
}

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
	case strings.HasPrefix(URL, "https://imhentai.xxx/"):
		return new(Parser_IMHENTAI_XXX), nil
	case strings.HasPrefix(URL, "https://www.3hentai1.buzz/"):
		return new(Parser_3HENTAI1_BUZZ), nil
	case strings.HasPrefix(URL, "https://manga-online.biz/"):
		return new(Parser_MANGAONLINE_BIZ), nil
	case strings.HasPrefix(URL, "https://doujins.com/"):
		return new(Parser_DOUJINS_COM), nil
	default:
		return nil, ErrInvalidLink
	}
}

func Load(ctx context.Context, r Requester, URL string) (p Parser, ok bool, err error) {
	p, err = Parse(ctx, URL)
	if err == nil {
		ok = p.Load(ctx, r, URL)
	}
	return p, ok, err
}

// Parser интерфейс для реализации парсеров для различных сайтов
type Parser interface {
	Load(ctx context.Context, r Requester, URL string) bool
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

func ParseAttr(ctx context.Context, p Parser, attr domain.Attribute) []string {
	switch attr {
	case domain.AttrAuthor:
		return p.ParseAuthors(ctx)

	case domain.AttrCategory:
		return p.ParseCategories(ctx)

	case domain.AttrCharacter:
		return p.ParseCharacters(ctx)

	case domain.AttrGroup:
		return p.ParseGroups(ctx)

	case domain.AttrLanguage:
		return p.ParseLanguages(ctx)

	case domain.AttrParody:
		return p.ParseParodies(ctx)

	case domain.AttrTag:
		return p.ParseTags(ctx)

	default:
		return []string{}
	}
}

// Проверка соответствия базового типа
var _ Parser = (*baseParser)(nil)

type baseParser struct{}

func (p baseParser) Load(ctx context.Context, r Requester, URL string) bool { return false }
func (p baseParser) ParseName(ctx context.Context) string                   { return "" }
func (p baseParser) ParsePages(ctx context.Context) []Page                  { return []Page{} }
func (p baseParser) ParseTags(ctx context.Context) []string                 { return []string{} }
func (p baseParser) ParseAuthors(ctx context.Context) []string              { return []string{} }
func (p baseParser) ParseCharacters(ctx context.Context) []string           { return []string{} }
func (p baseParser) ParseLanguages(ctx context.Context) []string            { return []string{} }
func (p baseParser) ParseCategories(ctx context.Context) []string           { return []string{} }
func (p baseParser) ParseParodies(ctx context.Context) []string             { return []string{} }
func (p baseParser) ParseGroups(ctx context.Context) []string               { return []string{} }
