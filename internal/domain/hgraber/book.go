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
