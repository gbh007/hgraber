package hgraber

import (
	"context"
	"time"
)

type BookInfoParsed struct {
	Name bool
	Page bool

	Attributes map[Attribute]bool
}

func (info BookInfoParsed) IsFullParsed(ctx context.Context) bool {
	for _, parsed := range info.Attributes {
		if !parsed {
			return false
		}
	}

	return info.Name && info.Page
}

type BookInfo struct {
	Parsed BookInfoParsed
	Name   string
	Rate   int

	Attributes map[Attribute][]string
}

type Book struct {
	ID      int
	Created time.Time
	URL     string

	Pages []Page
	Data  BookInfo
}

type BookFilter struct {
	Limit    int
	Offset   int
	NewFirst bool
}
