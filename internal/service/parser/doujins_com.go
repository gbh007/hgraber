package parser

import (
	"app/system"
	"context"
	"html"
	"net/url"
	"regexp"
	"strings"
)

// Parser_DOUJINS_COM парсер для сайта https://doujins.com/
type Parser_DOUJINS_COM struct {
	baseParser

	main_raw string
	url      string
}

func (p *Parser_DOUJINS_COM) Load(ctx context.Context, URL string) bool {
	var err error
	p.url = URL
	p.main_raw, err = RequestString(ctx, URL)
	return err == nil
}

func (p Parser_DOUJINS_COM) ParseName(ctx context.Context) string {
	rp := `(?sm)` + regexp.QuoteMeta(`<title>`) +
		`(.+?)` + regexp.QuoteMeta(`</title>`)

	res := regexp.MustCompile(rp).FindAllStringSubmatch(p.main_raw, -1)
	if len(res) < 1 || len(res[0]) != 2 {
		return ""
	}

	return strings.TrimSpace(html.UnescapeString(res[0][1]))
}

func (p Parser_DOUJINS_COM) ParseTags(ctx context.Context) []string {
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

	return result
}

func (p Parser_DOUJINS_COM) parsePages(s string) []string {
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

func (p Parser_DOUJINS_COM) ParsePages(ctx context.Context) []Page {
	result := make([]Page, 0)
	res := p.parsePages(p.main_raw)
	if len(res) < 1 {
		return []Page{}
	}

	for i, rURL := range res {
		u, err := url.Parse(rURL)
		if err != nil {
			system.Error(ctx, err)

			continue
		}

		fnameTmp := strings.Split(u.Path, "/")                   // название файла
		fnameTmp = strings.Split(fnameTmp[len(fnameTmp)-1], ".") // расширение
		result = append(result, Page{URL: rURL, Number: i + 1, Ext: fnameTmp[len(fnameTmp)-1]})
	}

	return result
}

func (p Parser_DOUJINS_COM) ParseAuthors(ctx context.Context) []string {
	artistBlock := `(?sm)` + regexp.QuoteMeta(`<div class="gallery-artist">`) +
		`(.+?)` + regexp.QuoteMeta(`</div>`)

	blockRes := regexp.MustCompile(artistBlock).FindAllStringSubmatch(p.main_raw, -1)
	if len(blockRes) < 1 || len(blockRes[0]) != 2 {
		return []string{}
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

	return result
}
