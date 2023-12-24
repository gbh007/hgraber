package parser

import (
	"app/internal/domain/hgraber"
	"context"
	"fmt"
	"regexp"
	"strings"
)

type mockParser struct {
	url, body string
	r         Requester
}

func (p *mockParser) Load(ctx context.Context, r Requester, URL string) error {
	p.r = r
	p.url = URL

	if len(p.url) > 1 && p.url[len(p.url)-1] == '/' {
		p.url = p.url[:len(p.url)-1]
	}

	body, err := r.RequestString(ctx, URL)
	if err != nil {
		return fmt.Errorf("mock parser: %w", err)
	}

	p.body = body

	return nil
}

func (p *mockParser) ParsePages(ctx context.Context) []hgraber.Page {
	result := make([]hgraber.Page, 0)

	rp := `(?sm)` + regexp.QuoteMeta(`<a href="`) + `(.+?)\.(.+?)` + regexp.QuoteMeta(`">`)
	for i, name := range regexp.MustCompile(rp).FindAllStringSubmatch(p.body, -1) {
		if len(name) > 1 {
			result = append(result, hgraber.Page{
				PageNumber: i + 1,
				URL:        p.url + "/" + strings.TrimSpace(name[1]) + "." + strings.TrimSpace(name[2]),
				Ext:        name[2],
			})
		}
	}

	return result
}

func (p *mockParser) ParseName(ctx context.Context) string {
	return "mock name"
}

func (p *mockParser) ParseTags(ctx context.Context) []string {
	return []string{"mock Tags"}
}

func (p *mockParser) ParseAuthors(ctx context.Context) []string {
	return []string{"mock Authors"}
}

func (p *mockParser) ParseCharacters(ctx context.Context) []string {
	return []string{"mock Characters"}
}

func (p *mockParser) ParseLanguages(ctx context.Context) []string {
	return []string{"mock Languages"}
}

func (p *mockParser) ParseCategories(ctx context.Context) []string {
	return []string{"mock Categories"}
}

func (p *mockParser) ParseParodies(ctx context.Context) []string {
	return []string{"mock Parodies"}
}

func (p *mockParser) ParseGroups(ctx context.Context) []string {
	return []string{"mock Groups"}
}
