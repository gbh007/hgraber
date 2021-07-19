package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parser_IMHENTAI_XXX парсер для сайта https://imhentai.xxx/
type Parser_IMHENTAI_XXX struct {
	main_raw string
	url      string
}

func (p *Parser_IMHENTAI_XXX) Load(URL string) bool {
	var err error
	p.url = URL
	p.main_raw, err = RequestString(URL)
	if err != nil {
		return false
	}
	return true
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

func (p Parser_IMHENTAI_XXX) ParseName() string {
	rp := `(?sm)` + regexp.QuoteMeta(`<div class="row gallery_first">`) + `.+?` +
		regexp.QuoteMeta(`<h1>`) + `(.+?)` + regexp.QuoteMeta(`</h1>`)
	res := regexp.MustCompile(rp).FindAllStringSubmatch(p.main_raw, -1)
	if len(res) < 1 || len(res[0]) != 2 {
		return ""
	}
	return res[0][1]
}
func (p Parser_IMHENTAI_XXX) ParseTags() []string {
	if strings.Index(p.main_raw, `<span class='tags_text'>Tags:</span>`) < 0 {
		return []string{}
	}
	raw := strings.Split(p.main_raw, `<span class='tags_text'>Tags:</span>`)[1]
	raw = strings.Split(raw, `</li>`)[0]
	return p.parseTags(raw)
}
func (p Parser_IMHENTAI_XXX) ParsePages() []Page {
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
	rp_img := regexp.MustCompile(`(?sm)` + regexp.QuoteMeta(`<img id="gimg" class="lazy`) + `.+?` + `src\=\"(.*?)\"`)
	for i := 1; i <= count; i++ {
		// символ / и так будет в конце
		data, err := RequestString(fmt.Sprintf("%s%d", u, i))
		if err != nil {
			return []Page{}
		}
		res := rp_img.FindStringSubmatch(data)
		if len(res) < 2 {
			return []Page{}
		}
		fnameTmp := strings.Split(res[1], "/")                   // название файла
		fnameTmp = strings.Split(fnameTmp[len(fnameTmp)-1], ".") // расширение
		result = append(result, Page{URL: res[1], Number: i, Ext: fnameTmp[len(fnameTmp)-1]})
	}
	return result
}
func (p Parser_IMHENTAI_XXX) ParseAuthors() []string {
	if strings.Index(p.main_raw, `<span class='tags_text'>Artists:</span>`) < 0 {
		return []string{}
	}
	raw := strings.Split(p.main_raw, `<span class='tags_text'>Artists:</span>`)[1]
	raw = strings.Split(raw, `</li>`)[0]
	return p.parseTags(raw)
}
func (p Parser_IMHENTAI_XXX) ParseCharacters() []string {
	if strings.Index(p.main_raw, `<span class='tags_text'>Characters:</span>`) < 0 {
		return []string{}
	}
	raw := strings.Split(p.main_raw, `<span class='tags_text'>Characters:</span>`)[1]
	raw = strings.Split(raw, `</li>`)[0]
	return p.parseTags(raw)
}
