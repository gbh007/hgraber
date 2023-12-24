package parser

import (
	"app/internal/domain/hgraber"
	"context"
	"regexp"
	"strconv"
	"strings"
)

// Parser_MANGAONLINE_BIZ парсер для сайта https://manga-online.biz/
type Parser_MANGAONLINE_BIZ struct {
	baseParser
	main_raw string
	url      string

	r Requester
}

func (p *Parser_MANGAONLINE_BIZ) Load(ctx context.Context, r Requester, URL string) error {
	p.r = r

	var err error

	p.url = URL
	tmpUrl := trimLastSlash(URL, 4) + ".html"
	p.main_raw, err = r.RequestString(ctx, tmpUrl)

	if err != nil {
		return err
	}

	return nil
}

func (p Parser_MANGAONLINE_BIZ) ParseName(ctx context.Context) string {
	rp := `(?sm)` + regexp.QuoteMeta(`<h1 class="header">`) + `\s*(.+?)\s*` + regexp.QuoteMeta(`</h1>`)
	res := regexp.MustCompile(rp).FindAllStringSubmatch(p.main_raw, -1)
	if len(res) < 1 || len(res[0]) != 2 {
		return ""
	}
	return res[0][1]
}

func (p Parser_MANGAONLINE_BIZ) ParsePages(ctx context.Context) []hgraber.Page {
	result := make([]hgraber.Page, 0)
	pcDataRaw, err := p.r.RequestString(ctx, p.url)
	if err != nil {
		return result
	}
	rp_img := regexp.MustCompile(regexp.QuoteMeta(`{"number":`) + `(\d+)` + regexp.QuoteMeta(`,"src":"`) + `(.*?)` + regexp.QuoteMeta(`","`))
	rp_base := regexp.MustCompile(regexp.QuoteMeta(`'srcBaseUrl': '`) + `(.*?)` + regexp.QuoteMeta(`',`))
	tmpBase := rp_base.FindStringSubmatch(pcDataRaw)
	if len(tmpBase) != 2 {
		return result
	}
	baseURL := tmpBase[1]
	for _, pg := range rp_img.FindAllStringSubmatch(pcDataRaw, -1) {
		i, err := strconv.Atoi(pg[1])
		if err != nil {
			return []hgraber.Page{}
		}
		res := baseURL + strings.ReplaceAll(pg[2], `\/`, `/`)
		fnameTmp := strings.Split(res, "/")                      // название файла
		fnameTmp = strings.Split(fnameTmp[len(fnameTmp)-1], ".") // расширение
		result = append(result, hgraber.Page{URL: res, PageNumber: i, Ext: fnameTmp[len(fnameTmp)-1]})
	}
	return result
}

func (p Parser_MANGAONLINE_BIZ) ParseTags(ctx context.Context) []string {
	result := make([]string, 0)
	rp := `(?sm)` + regexp.QuoteMeta(`<a onclick="App.Analytics.track('Genre', 'Click', 'Manga');" href="`) + `.+?` +
		regexp.QuoteMeta(`" target="_blank" class="ui label">`) + `(.+?)` + regexp.QuoteMeta(`</a>`)
	for _, tag := range regexp.MustCompile(rp).FindAllStringSubmatch(p.main_raw, -1) {
		if len(tag) > 1 {
			result = append(result, strings.TrimSpace(tag[1]))
		}
	}
	return result
}
