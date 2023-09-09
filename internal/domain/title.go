package domain

import (
	"context"
	"time"
)

type TitleInfoParsed struct {
	Name bool
	Page bool

	Attributes map[Attribute]bool
}

func (tip TitleInfoParsed) IsFullParsed(ctx context.Context) bool {
	for _, parsed := range tip.Attributes {
		if !parsed {
			return false
		}
	}

	return tip.Name && tip.Page
}

type TitleInfo struct {
	Parsed TitleInfoParsed
	Name   string
	Rate   int

	Attributes map[Attribute][]string
}

type Title struct {
	ID      int
	Created time.Time
	URL     string

	Pages []Page
	Data  TitleInfo
}

type BookFilter struct {
	Limit    int
	Offset   int
	NewFirst bool
}
