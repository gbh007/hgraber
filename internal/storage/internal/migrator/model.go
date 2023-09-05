package migrator

import (
	"time"
)

// Migration - модель для таблицы с миграциями
type Migration struct {
	ID       int       `db:"id"`
	Filename string    `db:"filename"`
	Hash     string    `db:"hash"`
	Applied  time.Time `db:"applied"`
}
