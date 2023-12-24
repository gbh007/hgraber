package jdb

import (
	"app/internal/domain/hgraber"
	"context"
)

func (db *Database) UpdateBookName(ctx context.Context, id int, name string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	title.Data.Name = name
	title.Data.Parsed.Name = true

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateAttributes(ctx context.Context, id int, attr hgraber.Attribute, data []string) error {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	title, ok := db.data.Titles[id]
	if !ok {
		return hgraber.BookNotFoundError
	}

	switch attr {
	case hgraber.AttrAuthor:
		title.Data.Authors = data
		title.Data.Parsed.Authors = true

	case hgraber.AttrCategory:
		title.Data.Categories = data
		title.Data.Parsed.Categories = true

	case hgraber.AttrCharacter:
		title.Data.Characters = data
		title.Data.Parsed.Characters = true

	case hgraber.AttrGroup:
		title.Data.Groups = data
		title.Data.Parsed.Groups = true

	case hgraber.AttrLanguage:
		title.Data.Languages = data
		title.Data.Parsed.Languages = true

	case hgraber.AttrParody:
		title.Data.Parodies = data
		title.Data.Parsed.Parodies = true

	case hgraber.AttrTag:
		title.Data.Tags = data
		title.Data.Parsed.Tags = true

	default:
		return hgraber.UnsupportedAttributeError
	}

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}
