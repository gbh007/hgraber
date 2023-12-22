package parser

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parser_IMHENTAI_XXX парсер для сайта https://imhentai.xxx/
type Parser_IMHENTAI_XXX struct {
	main_raw string
	url      string

	r Requester
}

func (p *Parser_IMHENTAI_XXX) Load(ctx context.Context, r Requester, URL string) bool {
	p.r = r

	var err error

	p.url = URL
	p.main_raw, err = r.RequestString(ctx, URL)

	return err == nil
}

// parseTags парсит теги авторов и тд
func (p Parser_IMHENTAI_XXX) parseTags(s string) []string {
	result := make([]string, 0)
	rp := `(?sm)` + regexp.QuoteMeta(`<a class='`) + `.+?` +
		regexp.QuoteMeta(`' href='`) + `.+?` +
		regexp.QuoteMeta(`'>`) + `(.+?)` +
		regexp.QuoteMeta(`<span class='badge'>`)
	for _, tag := range regexp.MustCompile(rp).FindAllStringSubmatch(s, -1) {
		if len(tag) > 1 {
			result = append(result, strings.TrimSpace(tag[1]))
		}
	}
	return result
}

func (p Parser_IMHENTAI_XXX) ParseName(ctx context.Context) string {
	rp := `(?sm)` + regexp.QuoteMeta(`<div class="row gallery_first">`) + `.+?` +
		regexp.QuoteMeta(`<h1>`) + `(.+?)` + regexp.QuoteMeta(`</h1>`)
	res := regexp.MustCompile(rp).FindAllStringSubmatch(p.main_raw, -1)
	if len(res) < 1 || len(res[0]) != 2 {
		return ""
	}
	return regexp.MustCompile(`<a.+?</a>`).ReplaceAllString(res[0][1], "")
	// return res[0][1]
}
func (p Parser_IMHENTAI_XXX) ParseTags(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `<span class='tags_text'>Tags:</span>`) {
		return []string{}
	}
	raw := strings.Split(p.main_raw, `<span class='tags_text'>Tags:</span>`)[1]
	raw = strings.Split(raw, `</li>`)[0]
	return p.parseTags(raw)
}
func (p Parser_IMHENTAI_XXX) ParsePages(ctx context.Context) []Page {
	result := make([]Page, 0)
	rp := `(?sm)` + regexp.QuoteMeta(`<li class="pages">Pages: `) + `(\d+).*?` + regexp.QuoteMeta(`</li>`)
	res := regexp.MustCompile(rp).FindStringSubmatch(p.main_raw)
	if len(res) < 2 {
		return []Page{}
	}
	count, err := strconv.Atoi(res[1])
	if err != nil {
		return []Page{}
	}
	u := strings.Replace(p.url, "gallery", "view", -1)
	rp_img := regexp.MustCompile(regexp.QuoteMeta(`<img id="gimg" class="lazy`) + `.+?` + `src="(.+?)" alt`)
	for i := 1; i <= count; i++ {
		// символ / и так будет в конце
		data, err := p.r.RequestString(ctx, fmt.Sprintf("%s%d", u, i))
		if err != nil {
			return []Page{}
		}
		res := rp_img.FindStringSubmatch(data)
		if len(res) < 2 {
			return []Page{}
		}
		url := res[1]
		if strings.Contains(url, "data-src=\"") {
			url = strings.Split(url, "data-src=\"")[1]
		}
		fnameTmp := strings.Split(url, "/")                      // название файла
		fnameTmp = strings.Split(fnameTmp[len(fnameTmp)-1], ".") // расширение
		result = append(result, Page{URL: url, Number: i, Ext: fnameTmp[len(fnameTmp)-1]})
	}
	return result
}
func (p Parser_IMHENTAI_XXX) ParseAuthors(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `<span class='tags_text'>Artists:</span>`) {
		return []string{}
	}
	raw := strings.Split(p.main_raw, `<span class='tags_text'>Artists:</span>`)[1]
	raw = strings.Split(raw, `</li>`)[0]
	return p.parseTags(raw)
}
func (p Parser_IMHENTAI_XXX) ParseCharacters(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `<span class='tags_text'>Characters:</span>`) {
		return []string{}
	}
	raw := strings.Split(p.main_raw, `<span class='tags_text'>Characters:</span>`)[1]
	raw = strings.Split(raw, `</li>`)[0]
	return p.parseTags(raw)
}

func (p Parser_IMHENTAI_XXX) ParseLanguages(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `<span class='tags_text'>Languages:</span>`) {
		return []string{}
	}
	raw := strings.Split(p.main_raw, `<span class='tags_text'>Languages:</span>`)[1]
	raw = strings.Split(raw, `</li>`)[0]
	return p.parseTags(raw)
}
func (p Parser_IMHENTAI_XXX) ParseCategories(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `<span class='tags_text'>Category:</span>`) {
		return []string{}
	}
	raw := strings.Split(p.main_raw, `<span class='tags_text'>Category:</span>`)[1]
	raw = strings.Split(raw, `</li>`)[0]
	return p.parseTags(raw)
}
func (p Parser_IMHENTAI_XXX) ParseParodies(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `<span class='tags_text'>Parodies:</span>`) {
		return []string{}
	}
	raw := strings.Split(p.main_raw, `<span class='tags_text'>Parodies:</span>`)[1]
	raw = strings.Split(raw, `</li>`)[0]
	return p.parseTags(raw)
}
func (p Parser_IMHENTAI_XXX) ParseGroups(ctx context.Context) []string {
	if !strings.Contains(p.main_raw, `<span class='tags_text'>Groups:</span>`) {
		return []string{}
	}
	raw := strings.Split(p.main_raw, `<span class='tags_text'>Groups:</span>`)[1]
	raw = strings.Split(raw, `</li>`)[0]
	return p.parseTags(raw)
}
