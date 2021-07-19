package db

import (
	"log"
	"time"
)

// type Title struct {
// 	ID               int
// 	Name             string
// 	URL              string
// 	PageCount        int
// 	CreationTime     time.Time
// 	Loaded           bool
// 	ParsedPages      bool
// 	ParsedTags       bool
// 	ParsedAuthors    bool
// 	ParsedCharacters bool
// }

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
