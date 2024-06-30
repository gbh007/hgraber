package hgraber

import (
	"time"
)

type BookInfoParsed struct {
	Name bool
	Page bool

	Attributes map[Attribute]bool
}

type BookInfo struct {
	Parsed BookInfoParsed
	Name   string
	Rating int

	PageCount int

	Attributes map[Attribute][]string
}

type Book struct {
	ID      int
	Created time.Time
	URL     string

	Pages []Page
	Data  BookInfo
}

func (b Book) PageCount() int {
	if b.Data.PageCount == 0 {
		return len(b.Pages)
	}

	return b.Data.PageCount
}

func (b Book) AttributesParsed() bool {
	for _, ok := range b.Data.Parsed.Attributes {
		if !ok {
			return false
		}
	}

	return true
}

type BookFilter struct {
	Limit    int
	Offset   int
	NewFirst bool
}
