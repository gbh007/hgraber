package db

import (
	"database/sql"
	"log"
)

// GetTagID возвращает ид тега, в случае его отсутствия создает
func GetTagID(name string) (int, error) {
	var id int
	row := _db.QueryRow(`SELECT id FROM tags WHERE name = ?`, name)
	err := row.Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		log.Println(err)
		return 0, err
	}
	result, err := _db.Exec(`INSERT INTO tags(name) VALUES(?)`, name)
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

// GetAuthorID возвращает ид автора, в случае его отсутствия создает
func GetAuthorID(name string) (int, error) {
	var id int
	row := _db.QueryRow(`SELECT id FROM authors WHERE name = ?`, name)
	err := row.Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		log.Println(err)
		return 0, err
	}
	result, err := _db.Exec(`INSERT INTO authors(name) VALUES(?)`, name)
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

// GetCharacterID возвращает ид персонажа, в случае его отсутствия создает
func GetCharacterID(name string) (int, error) {
	var id int
	row := _db.QueryRow(`SELECT id FROM characters WHERE name = ?`, name)
	err := row.Scan(&id)
	if err == nil {
		return id, nil
	}
	if err != sql.ErrNoRows {
		log.Println(err)
		return 0, err
	}
	result, err := _db.Exec(`INSERT INTO characters(name) VALUES(?)`, name)
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

// UpdateTitleTags обновляет в тайтле список тэгов
func UpdateTitleTags(id int, tags []string) error {
	ids := []int{}
	for _, tag := range tags {
		i, err := GetTagID(tag)
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
	_, err = tx.Exec(`DELETE FROM link_tags_titles WHERE title_id = ?`, id)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	// добавление новых данных
	for _, tid := range ids {
		_, err = tx.Exec(`INSERT INTO link_tags_titles(title_id, tag_id) VALUES(?, ?)`, id, tid)
		if err != nil {
			log.Println(err)
			_ = tx.Rollback()
			return err
		}
	}
	// обновление данных тайтла
	_, err = tx.Exec(`UPDATE titles SET parsed_tags = ? WHERE id = ?`, true, id)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// UpdateTitleAuthors обновляет в тайтле список авторов
func UpdateTitleAuthors(id int, authors []string) error {
	ids := []int{}
	for _, au := range authors {
		i, err := GetAuthorID(au)
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
	_, err = tx.Exec(`DELETE FROM link_authors_titles WHERE title_id = ?`, id)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	// добавление новых данных
	for _, tid := range ids {
		_, err = tx.Exec(`INSERT INTO link_authors_titles(title_id, author_id) VALUES(?, ?)`, id, tid)
		if err != nil {
			log.Println(err)
			_ = tx.Rollback()
			return err
		}
	}
	// обновление данных тайтла
	_, err = tx.Exec(`UPDATE titles SET parsed_authors = ? WHERE id = ?`, true, id)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

// UpdateTitleCharacters обновляет в тайтле список персонажей
func UpdateTitleCharacters(id int, authors []string) error {
	ids := []int{}
	for _, ch := range authors {
		i, err := GetCharacterID(ch)
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
	_, err = tx.Exec(`DELETE FROM link_characters_titles WHERE title_id = ?`, id)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	// добавление новых данных
	for _, tid := range ids {
		_, err = tx.Exec(`INSERT INTO link_characters_titles(title_id, character_id) VALUES(?, ?)`, id, tid)
		if err != nil {
			log.Println(err)
			_ = tx.Rollback()
			return err
		}
	}
	// обновление данных тайтла
	_, err = tx.Exec(`UPDATE titles SET parsed_characters = ? WHERE id = ?`, true, id)
	if err != nil {
		log.Println(err)
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
