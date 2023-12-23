package postgresql

import (
	"database/sql"
	"time"
)

type Book struct {
	ID        int            `db:"id"`
	Name      sql.NullString `db:"name"`
	Url       sql.NullString `db:"url"`
	PageCount sql.NullInt32  `db:"page_count"`
	CreateAt  time.Time      `db:"create_at"`
	Rate      int            `db:"rate"`
}

type Page struct {
	BookID     int          `db:"book_id"`
	PageNumber int          `db:"page_number"`
	Ext        string       `db:"ext"`
	Url        string       `db:"url"`
	Success    bool         `db:"success"`
	CreateAt   time.Time    `db:"create_at"`
	LoadAt     sql.NullTime `db:"load_at"`
	Rate       int          `db:"rate"`
}

type Attribute struct {
	Code string `db:"code"`
}

type BookAttribute struct {
	BookID int    `db:"book_id"`
	Attr   string `db:"attr"`
	Value  string `db:"value"`
}

type BookAttributeParsed struct {
	BookID int    `db:"book_id"`
	Attr   string `db:"attr"`
	Parsed bool   `db:"parsed"`
}
