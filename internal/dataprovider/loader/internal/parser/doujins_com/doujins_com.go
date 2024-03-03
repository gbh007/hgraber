package doujins_com

import (
	"app/internal/dataprovider/loader/internal/parser/common"
	"app/internal/domain/hgraber"
	"context"
	"errors"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
)

// Проверка соответствия базового типа
var (
	_ hgraber.BookParser = (*BookParser)(nil)
	_ hgraber.Parser     = (*Parser)(nil)

	ParserError = errors.New("parser doujins.com")
)

// Parser парсер для сайта https://doujins.com/
type Parser struct {
	common.CoreParser
}

func New(r common.Requester) *Parser {
	return &Parser{
		CoreParser: common.NewCoreParser(r, []string{
			"https://doujins.com/",
		}),
	}
}

func (p *Parser) Load(ctx context.Context, URL string) (hgraber.BookParser, error) {
	bookParser := BookParser{
		url: URL,
	}

	var err error

	bookParser.main_raw, err = p.Requester.RequestString(ctx, URL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ParserError, err)
	}

	return bookParser, nil
}

// BookParser парсер для сайта https://doujins.com/
type BookParser struct {
	main_raw string
	url      string
}

func (p BookParser) Name(ctx context.Context) (string, error) {
	rp := `(?sm)` + regexp.QuoteMeta(`<title>`) +
		`(.+?)` + regexp.QuoteMeta(`</title>`)

	res := regexp.MustCompile(rp).FindAllStringSubmatch(p.main_raw, -1)
	if len(res) < 1 || len(res[0]) != 2 {
		return "", fmt.Errorf("%w: missing name", ParserError)
	}

	return strings.TrimSpace(html.UnescapeString(res[0][1])), nil
}

func (p BookParser) Tags(ctx context.Context) ([]string, error) {
	tmp := make(map[string]struct{})

	result := make([]string, 0)

	rp := `(?sm)` + regexp.QuoteMeta(`<a href="/searches?tag_id=`) +
		`\d+` +
		regexp.QuoteMeta(`" class="">`) +
		`(.+?)` +
		regexp.QuoteMeta(`</a>`)

	for _, tags := range regexp.MustCompile(rp).FindAllStringSubmatch(p.main_raw, -1) {
		if len(tags) > 1 {
			tmp[strings.TrimSpace(tags[1])] = struct{}{}
		}
	}

	for tag := range tmp {
		result = append(result, tag)
	}

	return result, nil
}

func (p BookParser) parsePages(s string) []string {
	result := make([]string, 0)

	rp := `(?sm)` + regexp.QuoteMeta(`<img id="`) +
		`.+?` +
		regexp.QuoteMeta(`" data-src="`) +
		`(.+?)` +
		regexp.QuoteMeta(`" class="swiper-lazy"/>`)

	for _, pages := range regexp.MustCompile(rp).FindAllStringSubmatch(s, -1) {
		if len(pages) > 1 {
			result = append(result, html.UnescapeString(strings.TrimSpace(pages[1])))
		}
	}

	return result
}

func (p BookParser) Pages(ctx context.Context) ([]hgraber.Page, error) {
	result := make([]hgraber.Page, 0)
	res := p.parsePages(p.main_raw)
	if len(res) < 1 {
		return nil, fmt.Errorf("%w: missing pages", ParserError)
	}

	for i, rURL := range res {
		u, err := url.Parse(rURL)
		if err != nil {
			return nil, err
		}

		fnameTmp := strings.Split(u.Path, "/")                   // название файла
		fnameTmp = strings.Split(fnameTmp[len(fnameTmp)-1], ".") // расширение
		result = append(result, hgraber.Page{URL: rURL, PageNumber: i + 1, Ext: fnameTmp[len(fnameTmp)-1]})
	}

	return result, nil
}

func (p BookParser) Authors(ctx context.Context) ([]string, error) {
	artistBlock := `(?sm)` + regexp.QuoteMeta(`<div class="gallery-artist">`) +
		`(.+?)` + regexp.QuoteMeta(`</div>`)

	blockRes := regexp.MustCompile(artistBlock).FindAllStringSubmatch(p.main_raw, -1)
	if len(blockRes) < 1 || len(blockRes[0]) != 2 {
		return []string{}, nil
	}

	result := make([]string, 0)

	rp := `(?sm)` + regexp.QuoteMeta(`>`) +
		`(.+?)` +
		regexp.QuoteMeta(`</a>`)

	for _, pages := range regexp.MustCompile(rp).FindAllStringSubmatch(blockRes[0][1], -1) {
		if len(pages) > 1 {
			result = append(result, strings.TrimSpace(pages[1]))
		}
	}

	return result, nil
}

func (BookParser) Characters(ctx context.Context) ([]string, error) { return []string{}, nil }
func (BookParser) Languages(ctx context.Context) ([]string, error)  { return []string{}, nil }
func (BookParser) Categories(ctx context.Context) ([]string, error) { return []string{}, nil }
func (BookParser) Parodies(ctx context.Context) ([]string, error)   { return []string{}, nil }
func (BookParser) Groups(ctx context.Context) ([]string, error)     { return []string{}, nil }
