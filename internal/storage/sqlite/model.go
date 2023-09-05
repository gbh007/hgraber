package sqlite

import (
	"app/system"
	"context"
	"database/sql"
	"time"
)

type Book struct {
	ID        int            `db:"id"`
	Name      sql.NullString `db:"name"`
	Url       sql.NullString `db:"url"`
	PageCount sql.NullInt32  `db:"page_count"`
	CreateAt  string         `db:"create_at"`
	Rate      int            `db:"rate"`
}

type Page struct {
	BookID     int            `db:"book_id"`
	PageNumber int            `db:"page_number"`
	Ext        string         `db:"ext"`
	Url        string         `db:"url"`
	Success    bool           `db:"success"`
	CreateAt   string         `db:"create_at"`
	LoadAt     sql.NullString `db:"load_at"`
	Rate       int            `db:"rate"`
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

func timeToString(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format(time.RFC3339Nano)
}

func timeToSQLString(t time.Time) sql.NullString {
	s := timeToString(t)

	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func stringToTime(ctx context.Context, s string) time.Time {
	if s == "" {
		return time.Time{}
	}

	t, err := time.ParseInLocation(time.RFC3339Nano, s, time.UTC)
	system.IfErr(ctx, err)

	return t
}
