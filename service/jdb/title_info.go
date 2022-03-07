package jdb

import (
	"app/system"
	"context"
)

func (db *Database) UpdateTitleName(ctx context.Context, id int, name string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "UpdateTitleName")()

	title, ok := db.data.Titles[id]
	if !ok {
		return TitleIndexError
	}

	title.Data.Name = name
	title.Data.Parsed.Name = true

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateTitleAuthors(ctx context.Context, id int, authors []string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "UpdateTitleAuthors")()

	title, ok := db.data.Titles[id]
	if !ok {
		return TitleIndexError
	}

	title.Data.Authors = authors
	title.Data.Parsed.Authors = true

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateTitleTags(ctx context.Context, id int, tags []string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "UpdateTitleTags")()

	title, ok := db.data.Titles[id]
	if !ok {
		return TitleIndexError
	}

	title.Data.Tags = tags
	title.Data.Parsed.Tags = true

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateTitleCharacters(ctx context.Context, id int, characters []string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "UpdateTitleCharacters")()

	title, ok := db.data.Titles[id]
	if !ok {
		return TitleIndexError
	}

	title.Data.Characters = characters
	title.Data.Parsed.Characters = true

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateTitleCategories(ctx context.Context, id int, categories []string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "UpdateTitleCategories")()

	title, ok := db.data.Titles[id]
	if !ok {
		return TitleIndexError
	}

	title.Data.Categories = categories
	title.Data.Parsed.Categories = true

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateTitleGroups(ctx context.Context, id int, groups []string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "UpdateTitleGroups")()

	title, ok := db.data.Titles[id]
	if !ok {
		return TitleIndexError
	}

	title.Data.Groups = groups
	title.Data.Parsed.Groups = true

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateTitleLanguages(ctx context.Context, id int, languages []string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "UpdateTitleLanguages")()

	title, ok := db.data.Titles[id]
	if !ok {
		return TitleIndexError
	}

	title.Data.Languages = languages
	title.Data.Parsed.Languages = true

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}

func (db *Database) UpdateTitleParodies(ctx context.Context, id int, parodies []string) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	defer system.Stopwatch(ctx, "UpdateTitleParodies")()

	title, ok := db.data.Titles[id]
	if !ok {
		return TitleIndexError
	}

	title.Data.Parodies = parodies
	title.Data.Parsed.Parodies = true

	db.data.Titles[id] = title
	db.needSave = true

	return nil

}