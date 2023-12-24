package parser

import (
	"app/internal/domain/hgraber"
	"context"
	"fmt"
	"regexp"
	"strings"
)

// Parser_3HENTAI1_BUZZ парсер для сайта https://www.3hentai1.buzz/
type Parser_3HENTAI1_BUZZ struct {
	main_raw string
	url      string

	r Requester
}

func (p *Parser_3HENTAI1_BUZZ) Load(ctx context.Context, r Requester, URL string) error {
	p.r = r

	var err error

	p.url = URL
	p.main_raw, err = r.RequestString(ctx, URL)

	if err != nil {
		return err
	}

	return nil
}

// parseTags парсит теги авторов и тд
func (p Parser_3HENTAI1_BUZZ) parseTags(s, sep string) []string {
	result := make([]string, 0)
	rp := `(?sm)` + regexp.QuoteMeta(`<span class="filter-elem"><a class="name" href="`+sep+"/") + `.+?` +
		regexp.QuoteMeta(`>`) + `(.+?)` +
		regexp.QuoteMeta(`</a></span>`)
	for _, tag := range regexp.MustCompile(rp).FindAllStringSubmatch(s, -1) {
		if len(tag) > 1 {
			result = append(result, strings.TrimSpace(tag[1]))
		}
	}
	return result
}

func (p Parser_3HENTAI1_BUZZ) ParseName(ctx context.Context) string {
	rp := `(?sm)` + regexp.QuoteMeta(`<h1 class="text-left font-weight-bold">`) + `.+?` +
		regexp.QuoteMeta(`<span class="middle-title">`) + `(.+?)` + regexp.QuoteMeta(`</span>`)
	res := regexp.MustCompile(rp).FindAllStringSubmatch(p.main_raw, -1)
	if len(res) < 1 || len(res[0]) != 2 {
		return ""
	}
	return strings.TrimSpace(res[0][1])
}

func (p Parser_3HENTAI1_BUZZ) ParseTags(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `Tags:`) {
		return []string{}
	}
	return p.parseTags(p.main_raw, "?tags")
}

func (p Parser_3HENTAI1_BUZZ) parsePages(s string) []string {
	result := make([]string, 0)

	rp := `(?sm)` + regexp.QuoteMeta(`<div class="single-thumb"><a href="`) +
		`(.+?)` +
		regexp.QuoteMeta(`" rel="nofollow"><img src=`)
	for _, tag := range regexp.MustCompile(rp).FindAllStringSubmatch(s, -1) {
		if len(tag) > 1 {
			result = append(result, strings.TrimSpace(tag[1]))
		}
	}
	return result
}

func (p Parser_3HENTAI1_BUZZ) ParsePages(ctx context.Context) []hgraber.Page {
	result := make([]hgraber.Page, 0)
	res := p.parsePages(p.main_raw)
	if len(res) < 1 {
		return []hgraber.Page{}
	}

	rp_img := regexp.MustCompile(regexp.QuoteMeta(`<img src="`) + `(.+?)` + regexp.QuoteMeta(`"`))
	for i, rURL := range res {
		// символ / и так будет в конце
		data, err := p.r.RequestString(ctx, fmt.Sprintf("https://www.3hentai1.buzz/%s", rURL))
		if err != nil {
			return []hgraber.Page{}
		}
		res := rp_img.FindStringSubmatch(data)
		if len(res) < 2 {
			return []hgraber.Page{}
		}
		url := res[1]

		fnameTmp := strings.Split(url, "/")                      // название файла
		fnameTmp = strings.Split(fnameTmp[len(fnameTmp)-1], ".") // расширение
		result = append(result, hgraber.Page{URL: url, PageNumber: i + 1, Ext: fnameTmp[len(fnameTmp)-1]})
	}

	return result
}

func (p Parser_3HENTAI1_BUZZ) ParseAuthors(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `Artists:`) {
		return []string{}
	}
	return p.parseTags(p.main_raw, "?artists")
}

func (p Parser_3HENTAI1_BUZZ) ParseCharacters(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `Characters:`) {
		return []string{}
	}
	return p.parseTags(p.main_raw, "?characters")
}

func (p Parser_3HENTAI1_BUZZ) ParseLanguages(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `Languages:`) {
		return []string{}
	}
	return p.parseTags(p.main_raw, "?language")
}

func (p Parser_3HENTAI1_BUZZ) ParseCategories(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `Categories:`) {
		return []string{}
	}
	return p.parseTags(p.main_raw, "?category")
}

func (p Parser_3HENTAI1_BUZZ) ParseParodies(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `Series:`) {
		return []string{}
	}
	return p.parseTags(p.main_raw, "?series")
}

func (p Parser_3HENTAI1_BUZZ) ParseGroups(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `Groups:`) {
		return []string{}
	}
	return p.parseTags(p.main_raw, "?groups")
}
