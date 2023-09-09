package fileStorage

import (
	"time"
)

const (
	interval = time.Second * 15
	// queueSize максимальный размер очереди для загрузки файлов страницы
	queueSize = 10000
	// handlersCount количество одновременно запущенных загрузчиков страниц
	handlersCount = 10
)

type qPage struct {
	TitleID    int
	PageNumber int
	URL        string
	Ext        string
}
