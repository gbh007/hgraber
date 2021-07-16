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
