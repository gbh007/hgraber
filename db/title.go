package db

import (
	"app/system/clog"
	"app/system/coreContext"
	"database/sql"
	"time"
)

type Page struct {
	TitleID    int    `json:"title_id"`
	PageNumber int    `json:"page_number"`
	URL        string `json:"url"`
	Ext        string `json:"ext"`
}

// InsertTitle добавляет тайтл
func InsertTitle(ctx coreContext.CoreContext, name, URL string, loaded bool) (int, error) {
	result, err := _db.ExecContext(
		ctx,
		`INSERT INTO titles(name, url, creation_time, loaded) VALUES(?, ?, ?, ?)`,
		name, URL, time.Now(), loaded,
	)
	if err != nil {
		clog.Error(ctx, err)
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		clog.Error(ctx, err)
		return -1, err
	}
	return int(id), nil
}

// UpdateTitle обновляет тайтл
func UpdateTitle(ctx coreContext.CoreContext, id int, name string, loaded bool) error {
	_, err := _db.ExecContext(
		ctx,
		`UPDATE titles SET name = ?, loaded = ? WHERE id = ?`,
		name, loaded, id,
	)
	if err != nil {
		clog.Error(ctx, err)
	}
	return err
}

// UpdateTitleParsedPage обновляет информацию об обработанных страницах в тайтле
func UpdateTitleParsedPage(ctx coreContext.CoreContext, id, count int, success bool) error {
	_, err := _db.ExecContext(ctx, `UPDATE titles SET parsed_pages = ?, page_count = ? WHERE id = ?`, success, count, id)
	if err != nil {
		clog.Error(ctx, err)
	}
	return err
}

// InsertPage добавляет страницу тайтла
func InsertPage(ctx coreContext.CoreContext, id int, ext, URL string, page_number int) error {
	_, err := _db.ExecContext(
		ctx,
		`INSERT INTO pages(title_id, ext, url, page_number, success) VALUES(?, ?, ?, ?, ?)
		ON CONFLICT(title_id, page_number) DO UPDATE SET ext = excluded.ext, url = excluded.url, success = false`,
		id, ext, URL, page_number, false,
	)
	if err != nil {
		clog.Error(ctx, err)
	}
	return err
}

// UpdatePageSuccess обновляет информацию об успешной загрузке страницы
func UpdatePageSuccess(ctx coreContext.CoreContext, id, page int, success bool) error {
	_, err := _db.ExecContext(ctx, `UPDATE pages SET success = ? WHERE title_id = ? AND page_number = ?`, success, id, page)
	if err != nil {
		clog.Error(ctx, err)
	}
	return err
}

// SelectUnsuccessPages выбирает из базы не загруженные страницы
func SelectUnsuccessPages(ctx coreContext.CoreContext) []Page {
	result := []Page{}
	rows, err := _db.QueryContext(ctx, `SELECT p.title_id, p.page_number, p.url, p.ext FROM
titles t INNER JOIN pages p ON t.parsed_pages = TRUE AND t.id = p.title_id AND p.success = FALSE`)
	if err != nil {
		clog.Error(ctx, err)
		return result
	}
	for rows.Next() {
		p := Page{}
		err = rows.Scan(&p.TitleID, &p.PageNumber, &p.URL, &p.Ext)
		if err != nil {
			clog.Error(ctx, err)
		} else {
			result = append(result, p)
		}
	}
	return result
}

// TitleShortInfo информация о тайтле
type TitleShortInfo struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	PageCount        int       `json:"page_count"`
	Loaded           bool      `json:"loaded"`
	ParsedPage       bool      `json:"parsed_page"`
	ParsedTags       bool      `json:"parsed_tags"`
	ParsedAuthors    bool      `json:"parsed_authors"`
	ParsedCharacters bool      `json:"parsed_characters"`
	Avg              float64   `json:"avg"`
	Ext              string    `json:"ext"`
	URL              string    `json:"url"`
	Created          time.Time `json:"created"`
	Tags             []string  `json:"tags"`
	Authors          []string  `json:"authors"`
	Characters       []string  `json:"characters"`
	ParsedLanguages  bool      `json:"parsed_languages"`
	ParsedCategories bool      `json:"parsed_categories"`
	ParsedParodies   bool      `json:"parsed_parodies"`
	ParsedGroups     bool      `json:"parsed_groups"`
	Languages        []string  `json:"languages"`
	Categories       []string  `json:"categories"`
	Parodies         []string  `json:"parodies"`
	Groups           []string  `json:"groups"`
}

// SelectTitles выбирает из базы все тайтлы
func SelectTitles(ctx coreContext.CoreContext, offset, limit int) []TitleShortInfo {
	result := []TitleShortInfo{}
	rows, err := _db.QueryContext(ctx, `SELECT
	t2.id,
	t2.name,
	t2.page_count,
	t2.loaded,
	t2.parsed_pages,
	t2.parsed_tags,
	t2.parsed_authors,
	t2.parsed_characters,
	t2.parsed_languages,
	t2.parsed_categories,
	t2.parsed_parodies,
	t2.parsed_groups,
	a.av,
	p2.ext,
	t2.url,
	t2.creation_time
FROM
	titles t2
LEFT JOIN
(
	SELECT
		p.title_id AS id,
		AVG(p.success) AS av
	FROM
		titles t
	INNER JOIN pages p ON
		t.id = p.title_id
	GROUP BY
		p.title_id) a ON
	t2.id = a.id
LEFT JOIN pages p2 ON
	p2.title_id = t2.id
	AND p2.page_number = 1
ORDER BY
	t2.id DESC
LIMIT ?, ?`, offset, limit)
	if err != nil {
		clog.Error(ctx, err)
		return result
	}
	for rows.Next() {
		t := TitleShortInfo{}
		avg := sql.NullFloat64{}
		ext := sql.NullString{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.PageCount,
			&t.Loaded,
			&t.ParsedPage,
			&t.ParsedTags,
			&t.ParsedAuthors,
			&t.ParsedCharacters,
			&t.ParsedLanguages,
			&t.ParsedCategories,
			&t.ParsedParodies,
			&t.ParsedGroups,
			&avg,
			&ext,
			&t.URL,
			&t.Created,
		)
		if err != nil {
			clog.Error(ctx, err)
		} else {
			t.Avg = avg.Float64 * 100
			t.Ext = ext.String
			t.Tags = SelectMetaByTitleIDAndType(ctx, t.ID, TagsMetaType)
			t.Authors = SelectMetaByTitleIDAndType(ctx, t.ID, AuthorsMetaType)
			t.Characters = SelectMetaByTitleIDAndType(ctx, t.ID, CharactersMetaType)
			t.Languages = SelectMetaByTitleIDAndType(ctx, t.ID, LanguagesMetaType)
			t.Categories = SelectMetaByTitleIDAndType(ctx, t.ID, CategoriesMetaType)
			t.Parodies = SelectMetaByTitleIDAndType(ctx, t.ID, ParodiesMetaType)
			t.Groups = SelectMetaByTitleIDAndType(ctx, t.ID, GroupsMetaType)
			result = append(result, t)
		}
	}
	return result
}

// SelectTitles выбирает из базы тайтл по ID
func SelectTitleByID(ctx coreContext.CoreContext, id int) (TitleShortInfo, error) {
	row := _db.QueryRowContext(ctx, `SELECT
	t2.id,
	t2.name,
	t2.page_count,
	t2.loaded,
	t2.parsed_pages,
	t2.parsed_tags,
	t2.parsed_authors,
	t2.parsed_characters,
	t2.parsed_languages,
	t2.parsed_categories,
	t2.parsed_parodies,
	t2.parsed_groups,
	a.av,
	p2.ext,
	t2.url,
	t2.creation_time
FROM
	titles t2
LEFT JOIN
(
	SELECT
		p.title_id AS id,
		AVG(p.success) AS av
	FROM
		titles t
	INNER JOIN pages p ON
		t.id = p.title_id
	GROUP BY
		p.title_id) a ON
	t2.id = a.id
LEFT JOIN pages p2 ON
	p2.title_id = t2.id
	AND p2.page_number = 1
WHERE t2.id = ?
ORDER BY
	t2.id DESC`, id)
	t := TitleShortInfo{}
	avg := sql.NullFloat64{}
	ext := sql.NullString{}
	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.PageCount,
		&t.Loaded,
		&t.ParsedPage,
		&t.ParsedTags,
		&t.ParsedAuthors,
		&t.ParsedCharacters,
		&t.ParsedLanguages,
		&t.ParsedCategories,
		&t.ParsedParodies,
		&t.ParsedGroups,
		&avg,
		&ext,
		&t.URL,
		&t.Created,
	)
	if err != nil {
		clog.Error(ctx, err)
	} else {
		t.Avg = avg.Float64 * 100
		t.Ext = ext.String
		t.Tags = SelectMetaByTitleIDAndType(ctx, t.ID, TagsMetaType)
		t.Authors = SelectMetaByTitleIDAndType(ctx, t.ID, AuthorsMetaType)
		t.Characters = SelectMetaByTitleIDAndType(ctx, t.ID, CharactersMetaType)
		t.Languages = SelectMetaByTitleIDAndType(ctx, t.ID, LanguagesMetaType)
		t.Categories = SelectMetaByTitleIDAndType(ctx, t.ID, CategoriesMetaType)
		t.Parodies = SelectMetaByTitleIDAndType(ctx, t.ID, ParodiesMetaType)
		t.Groups = SelectMetaByTitleIDAndType(ctx, t.ID, GroupsMetaType)
	}
	return t, err
}

// SelectPagesByTitleID выбирает из базы все страницы из тайтла
func SelectPagesByTitleID(ctx coreContext.CoreContext, id int) []Page {
	result := []Page{}
	rows, err := _db.QueryContext(ctx, `SELECT p.title_id, p.page_number, p.url, p.ext
FROM pages p
WHERE p.success = TRUE AND p.title_id = ?
ORDER BY p.page_number`, id)
	if err != nil {
		clog.Error(ctx, err)
		return result
	}
	for rows.Next() {
		p := Page{}
		err = rows.Scan(
			&p.TitleID,
			&p.PageNumber,
			&p.URL,
			&p.Ext,
		)
		if err != nil {
			clog.Error(ctx, err)
		} else {
			result = append(result, p)
		}
	}
	return result
}

// SelectPagesByTitleIDAndNumber выбирает из базы все страницу из тайтла по его ид и номеру страницы
func SelectPagesByTitleIDAndNumber(ctx coreContext.CoreContext, id, pageNumber int) (Page, error) {
	row := _db.QueryRowContext(ctx, `SELECT p.title_id, p.page_number, p.url, p.ext
FROM pages p
WHERE p.success = TRUE AND p.title_id = ? AND p.page_number = ?`, id, pageNumber)
	p := Page{}
	err := row.Scan(
		&p.TitleID,
		&p.PageNumber,
		&p.URL,
		&p.Ext,
	)
	if err != nil {
		clog.Error(ctx, err)
	}
	return p, err
}

// SelectUnloadTitles выбирает из базы все недогруженые тайтлы
func SelectUnloadTitles(ctx coreContext.CoreContext) []TitleShortInfo {
	result := []TitleShortInfo{}
	rows, err := _db.QueryContext(ctx, `SELECT id, url, loaded, parsed_pages, parsed_tags, parsed_authors, parsed_characters,
	parsed_languages, parsed_categories, parsed_parodies, parsed_groups
	FROM titles WHERE loaded = FALSE OR parsed_pages = FALSE OR parsed_tags = FALSE OR parsed_authors = FALSE OR parsed_characters = FALSE
	OR parsed_languages = FALSE OR parsed_categories = FALSE OR parsed_parodies = FALSE OR parsed_groups = FALSE`)
	if err != nil {
		clog.Error(ctx, err)
		return result
	}
	for rows.Next() {
		t := TitleShortInfo{}
		err = rows.Scan(
			&t.ID,
			&t.URL,
			&t.Loaded,
			&t.ParsedPage,
			&t.ParsedTags,
			&t.ParsedAuthors,
			&t.ParsedCharacters,
			&t.ParsedLanguages,
			&t.ParsedCategories,
			&t.ParsedParodies,
			&t.ParsedGroups,
		)
		if err != nil {
			clog.Error(ctx, err)
		} else {
			result = append(result, t)
		}
	}
	return result
}

// SelectTitlesCount получает количество тайтлов в базе
func SelectTitlesCount(ctx coreContext.CoreContext) int {
	row := _db.QueryRowContext(ctx, `SELECT COUNT(id) FROM titles`)
	var c int
	err := row.Scan(&c)
	if err != nil {
		clog.Error(ctx, err)
	}
	return c
}

// SelectUnloadTitlesCount получает количество недогруженных тайтлов в базе
func SelectUnloadTitlesCount(ctx coreContext.CoreContext) int {
	row := _db.QueryRowContext(ctx, `SELECT COUNT(id) FROM titles WHERE 
	loaded = FALSE OR parsed_pages = FALSE OR parsed_tags = FALSE OR parsed_authors = FALSE OR parsed_characters = FALSE
	OR parsed_languages = FALSE OR parsed_categories = FALSE OR parsed_parodies = FALSE OR parsed_groups = FALSE
	`)
	var c int
	err := row.Scan(&c)
	if err != nil {
		clog.Error(ctx, err)
	}
	return c
}

// SelectPagesCount получает количество страниц в базе
func SelectPagesCount(ctx coreContext.CoreContext) int {
	row := _db.QueryRowContext(ctx, `SELECT COUNT(url) FROM pages`)
	var c int
	err := row.Scan(&c)
	if err != nil {
		clog.Error(ctx, err)
	}
	return c
}

// SelectUnloadPagesCount получает количество недогруженных страниц в базе
func SelectUnloadPagesCount(ctx coreContext.CoreContext) int {
	row := _db.QueryRowContext(ctx, `SELECT COUNT(url) FROM pages WHERE success = FALSE`)
	var c int
	err := row.Scan(&c)
	if err != nil {
		clog.Error(ctx, err)
	}
	return c
}
