package jdb

import (
	"app/internal/domain"
	"context"
)

func (db *Database) UpdateBookName(ctx context.Context, id int, name string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return domain.BookNotFoundError
	}

	title.Data.Name = name
	title.Data.Parsed.Name = true

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateAttributes(ctx context.Context, id int, attr domain.Attribute, data []string) error {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return domain.BookNotFoundError
	}

	switch attr {
	case domain.AttrAuthor:
		title.Data.Authors = data
		title.Data.Parsed.Authors = true

	case domain.AttrCategory:
		title.Data.Categories = data
		title.Data.Parsed.Categories = true

	case domain.AttrCharacter:
		title.Data.Characters = data
		title.Data.Parsed.Characters = true

	case domain.AttrGroup:
		title.Data.Groups = data
		title.Data.Parsed.Groups = true

	case domain.AttrLanguage:
		title.Data.Languages = data
		title.Data.Parsed.Languages = true

	case domain.AttrParody:
		title.Data.Parodies = data
		title.Data.Parsed.Parodies = true

	case domain.AttrTag:
		title.Data.Tags = data
		title.Data.Parsed.Tags = true

	default:
		return domain.UnsupportedAttributeError
	}

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}
