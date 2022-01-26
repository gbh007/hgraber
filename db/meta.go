package db

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	TagsMetaType       = "tags"
	AuthorsMetaType    = "authors"
	CharactersMetaType = "characters"
	LanguagesMetaType  = "languages"
	CategoriesMetaType = "categories"
	ParodiesMetaType   = "parodies"
	GroupsMetaType     = "groups"
)

// GetMetaID возвращает ид меты, в случае его отсутствия создает
func GetMetaID(name, tp string) (int, error) {
	var id int
	row := _db.QueryRow(`SELECT id FROM meta WHERE name = ? AND type = ?`, name, tp)
	err := row.Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		log.Println(err)
		return 0, err
	}
	result, err := _db.Exec(`INSERT INTO meta(name, type) VALUES(?, ?)`, name, tp)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	id64, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return int(id64), nil
}

// UpdateTitleMeta обновляет в тайтле список меты
func UpdateTitleMeta(id int, tp string, names []string) error {
	ids := []int{}
	for _, name := range names {
		i, err := GetMetaID(name, tp)
		if err != nil {
			log.Println(err)
			return err
		}
		ids = append(ids, i)
	}
	tx, err := _db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	// удаление старых данных
	_, err = tx.Exec(`DELETE FROM link_meta_titles WHERE title_id = ? AND type = ?`, id, tp)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	// добавление новых данных
	for _, mid := range ids {
		_, err = tx.Exec(`INSERT INTO link_meta_titles(title_id, meta_id, type) VALUES(?, ?, ?)`, id, mid, tp)
		if err != nil {
			log.Println(err)
			_ = tx.Rollback()
			return err
		}
	}
	// обновление данных тайтла
	_, err = tx.Exec(fmt.Sprintf("UPDATE titles SET parsed_%s = ? WHERE id = ?", tp), true, id)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// SelectMetaByTitleIDAndType получает мету тайтла по его ID и типу
func SelectMetaByTitleIDAndType(id int, tp string) []string {
	result := []string{}
	rows, err := _db.Query(`SELECT m.name
FROM link_meta_titles lmt INNER JOIN meta m ON lmt.meta_id = m.id 
WHERE lmt.title_id = ? AND lmt.type = ?
ORDER BY name`, id, tp)
	if err != nil {
		log.Println(err)
		return result
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Println(err)
		} else {
			result = append(result, name)
		}
	}
	return result
}