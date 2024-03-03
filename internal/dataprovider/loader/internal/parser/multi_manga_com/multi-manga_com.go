package multi_manga_com

import (
	"app/internal/dataprovider/loader/internal/parser/common"
	"app/internal/domain/hgraber"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Проверка соответствия базового типа
var (
	_ hgraber.BookParser = (*BookParser)(nil)
	_ hgraber.Parser     = (*Parser)(nil)

	ParserError = errors.New("parser multi-manga.com")
)

type Parser struct {
	common.CoreParser
}

func New(r common.Requester) *Parser {
	return &Parser{
		CoreParser: common.NewCoreParser(r, []string{
			"https://ww.multi-manga.com/",
			"https://w2.multi-manga.com/",
			"https://w3.multi-manga.com/",
		}),
	}
}

func (p *Parser) Load(ctx context.Context, URL string) (hgraber.BookParser, error) {
	bookParser := BookParser{
		r:   p.Requester,
		url: URL,
	}

	var err error

	bookParser.body, err = p.Requester.RequestString(ctx, URL)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ParserError, err)
	}

	return bookParser, nil
}

type BookParser struct {
	r common.Requester

	body string
	url  string
}

func (p BookParser) Name(ctx context.Context) (string, error) {
	rp := `(?sm)` + regexp.QuoteMeta(`<div id="info">`) + `\s*` +
		regexp.QuoteMeta(`<h1>`) +
		`(.+?)` + regexp.QuoteMeta(`</h1>`)

	res, ok := common.OneMatch(regexp.MustCompile(rp), p.body)
	if !ok {
		return "", fmt.Errorf("%w: missing name", ParserError)
	}

	return res, nil
}

func (p BookParser) Pages(ctx context.Context) ([]hgraber.Page, error) {
	pageCountRgx := `(?ism)` + regexp.QuoteMeta(`<div>`) +
		`(\d+?)` + regexp.QuoteMeta(` страниц(ы)</div>`)

	res, ok := common.OneMatch(regexp.MustCompile(pageCountRgx), p.body)
	if !ok {
		return nil, fmt.Errorf("%w: missing pages", ParserError)
	}

	pageCount, err := strconv.Atoi(res)
	if err != nil {
		return nil, fmt.Errorf("%w: page count: %w", ParserError, err)
	}

	result := make([]hgraber.Page, pageCount)

	imageRgx := regexp.MustCompile(`(?sm)` + regexp.QuoteMeta(`<img class="fit-horizontal" src="`) +
		`(.+?)` + regexp.QuoteMeta(`"`))

	for i := 1; i <= pageCount; i++ {
		urlToParse := p.url + "/" + strconv.Itoa(i) + "/"
		body, err := p.r.RequestString(ctx, urlToParse)
		if err != nil {
			return nil, fmt.Errorf("%w: page: %w", ParserError, err)
		}

		res, ok := common.OneMatch(imageRgx, body)
		if !ok {
			return nil, fmt.Errorf("%w: missing match page", ParserError)
		}

		imageUrl := strings.TrimSpace(res)
		extRaw := strings.Split(imageUrl, ".")
		if len(extRaw) < 2 {
			return nil, fmt.Errorf("%w: missing page extension", ParserError)
		}

		result[i-1] = hgraber.Page{
			PageNumber: i,
			URL:        imageUrl,
			Ext:        extRaw[len(extRaw)-1],
		}
	}

	return result, nil
}

func (p BookParser) parseTags(name string) []string {
	rgx := regexp.MustCompile(`(?ism)` + regexp.QuoteMeta(`<div class="tag-container field-name "> `) + name + regexp.QuoteMeta(` <span class="tags">`) +
		`(.+?)` + regexp.QuoteMeta(`</span>`))

	res, ok := common.OneMatch(rgx, p.body)
	if !ok {
		return nil
	}

	tagRgx := regexp.MustCompile(`(?sm)` + regexp.QuoteMeta(`>`) +
		`(.+?)` + regexp.QuoteMeta(`</a>`))

	out := make([]string, 0, 10)

	for _, tag := range tagRgx.FindAllStringSubmatch(res, -1) {
		if len(tag) > 1 {
			out = append(out, tag[1])
		}
	}

	return out
}

func (p BookParser) Tags(_ context.Context) ([]string, error) {
	return p.parseTags("Теги"), nil
}

func (p BookParser) Authors(_ context.Context) ([]string, error) {
	return p.parseTags("Автор"), nil
}

func (p BookParser) Languages(_ context.Context) ([]string, error) {
	return p.parseTags("Язык"), nil
}

func (p BookParser) Parodies(_ context.Context) ([]string, error) {
	rgx := regexp.MustCompile(`(?ism)` + regexp.QuoteMeta(`<div class="tag-container field-name "> Серия <span class="tags">`) +
		`(.+?)` + regexp.QuoteMeta(`</span>`))

	res, ok := common.OneMatch(rgx, p.body)
	if !ok {
		return []string{}, nil
	}

	return []string{res}, nil
}

func (BookParser) Characters(_ context.Context) ([]string, error) { return []string{}, nil }
func (BookParser) Categories(_ context.Context) ([]string, error) { return []string{}, nil }
func (BookParser) Groups(_ context.Context) ([]string, error)     { return []string{}, nil }
