package hgraber

import (
	"context"
	"errors"
)

var ErrInvalidLink = errors.New("invalid link")

type Parser interface {
	ParseName(ctx context.Context) string
	ParsePages(ctx context.Context) []Page
	ParseTags(ctx context.Context) []string
	ParseAuthors(ctx context.Context) []string
	ParseCharacters(ctx context.Context) []string
	ParseLanguages(ctx context.Context) []string
	ParseCategories(ctx context.Context) []string
	ParseParodies(ctx context.Context) []string
	ParseGroups(ctx context.Context) []string
}

func ParseAttr(ctx context.Context, p Parser, attr Attribute) []string {
	switch attr {
	case AttrAuthor:
		return p.ParseAuthors(ctx)

	case AttrCategory:
		return p.ParseCategories(ctx)

	case AttrCharacter:
		return p.ParseCharacters(ctx)

	case AttrGroup:
		return p.ParseGroups(ctx)

	case AttrLanguage:
		return p.ParseLanguages(ctx)

	case AttrParody:
		return p.ParseParodies(ctx)

	case AttrTag:
		return p.ParseTags(ctx)

	default:
		return []string{}
	}
}
