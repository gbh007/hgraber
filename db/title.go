package db

import (
	"database/sql"
	"log"
	"time"
)

type Page struct {
	TitleID    int
	PageNumber int
	URL        string
	Ext        string
}

// InsertTitle добавляет тайтл
func InsertTitle(name, URL string, loaded bool) (int, error) {
	result, err := _db.Exec(
		`INSERT INTO titles(name, url, creation_time, loaded) VALUES(?, ?, ?, ?)`,
		name, URL, time.Now(), loaded,
	)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return int(id), nil
}

// UpdateTitle обновляет тайтл
func UpdateTitle(id int, name string, loaded bool) error {
	_, err := _db.Exec(
		`UPDATE titles SET name = ?, loaded = ? WHERE id = ?`,
		name, loaded, id,
	)
	if err != nil {
		log.Println(err)
	}
	return err
}

// UpdateTitleParsedPage обновляет информацию об обработанных страницах в тайтле
func UpdateTitleParsedPage(id, count int, success bool) error {
	_, err := _db.Exec(`UPDATE titles SET parsed_pages = ?, page_count = ? WHERE id = ?`, success, count, id)
	if err != nil {
		log.Println(err)
	}
	return err
}

// InsertPage добавляет страницу тайтла
func InsertPage(id int, ext, URL string, page_number int) error {
	_, err := _db.Exec(
		`INSERT INTO pages(title_id, ext, url, page_number, success) VALUES(?, ?, ?, ?, ?)
		ON CONFLICT(title_id, page_number) DO UPDATE SET ext = excluded.ext, url = excluded.url, success = false`,
		id, ext, URL, page_number, false,
	)
	if err != nil {
		log.Println(err)
	}
	return err
}

// UpdatePageSuccess обновляет информацию об успешной загрузке страницы
func UpdatePageSuccess(id, page int, success bool) error {
	_, err := _db.Exec(`UPDATE pages SET success = ? WHERE title_id = ? AND page_number = ?`, success, id, page)
	if err != nil {
		log.Println(err)
	}
	return err
}

// SelectUnsuccessPages выбирает из базы не загруженные страницы
func SelectUnsuccessPages() []Page {
	result := []Page{}
	rows, err := _db.Query(`SELECT p.title_id, p.page_number, p.url, p.ext FROM
titles t INNER JOIN pages p ON t.loaded = TRUE AND t.parsed_pages = TRUE AND t.id = p.title_id AND p.success = FALSE`)
	if err != nil {
		log.Println(err)
		return result
	}
	for rows.Next() {
		p := Page{}
		err = rows.Scan(&p.TitleID, &p.PageNumber, &p.URL, &p.Ext)
		if err != nil {
			log.Println(err)
		} else {
			result = append(result, p)
		}
	}
	return result
}

// TitleShortInfo информация о тайтле
type TitleShortInfo struct {
	ID         int     `json:"id"`
	Name       string  `json:"name"`
	PageCount  int     `json:"pc"`
	Loaded     bool    `json:"loaded"`
	ParsedPage bool    `json:"pp"`
	Avg        float64 `json:"avg"`
	Ext        string  `json:"ext"`
	URL        string  `json:"url"`
}

// SelectTitles выбирает из базы все тайтлы
func SelectTitles() []TitleShortInfo {
	result := []TitleShortInfo{}
	rows, err := _db.Query(`SELECT
	t2.id,
	t2.name,
	t2.page_count,
	t2.loaded,
	t2.parsed_pages,
	a.av,
	p2.ext,
	t2.url
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
	t2.id DESC`)
	if err != nil {
		log.Println(err)
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
			&avg,
			&ext,
			&t.URL,
		)
		if err != nil {
			log.Println(err)
		} else {
			t.Avg = avg.Float64 * 100
			t.Ext = ext.String
			result = append(result, t)
		}
	}
	return result
}

// SelectTitles выбирает из базы тайтл по ID
func SelectTitleByID(id int) (TitleShortInfo, error) {
	row := _db.QueryRow(`SELECT
	t2.id,
	t2.name,
	t2.page_count,
	t2.loaded,
	t2.parsed_pages,
	a.av,
	p2.ext,
	t2.url
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
		&avg,
		&ext,
		&t.URL,
	)
	if err != nil {
		log.Println(err)
	} else {
		t.Avg = avg.Float64 * 100
		t.Ext = ext.String
	}
	return t, err
}

// SelectPagesByTitleID выбирает из базы все страницы из тайтла
func SelectPagesByTitleID(id int) []Page {
	result := []Page{}
	rows, err := _db.Query(`SELECT p.title_id, p.page_number, p.url, p.ext
FROM pages p
WHERE p.success = TRUE AND p.title_id = ?
ORDER BY p.page_number`, id)
	if err != nil {
		log.Println(err)
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
			log.Println(err)
		} else {
			result = append(result, p)
		}
	}
	return result
}

// SelectUnloadTitles выбирает из базы все недогруженые тайтлы
func SelectUnloadTitles() []TitleShortInfo {
	result := []TitleShortInfo{}
	rows, err := _db.Query(`SELECT t.id, t.url FROM	titles t WHERE t.loaded = FALSE OR t.parsed_pages = FALSE`)
	if err != nil {
		log.Println(err)
		return result
	}
	for rows.Next() {
		t := TitleShortInfo{}
		err = rows.Scan(
			&t.ID,
			&t.URL,
		)
		if err != nil {
			log.Println(err)
		} else {
			result = append(result, t)
		}
	}
	return result
}
