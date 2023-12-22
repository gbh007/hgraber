package pageHandler

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
