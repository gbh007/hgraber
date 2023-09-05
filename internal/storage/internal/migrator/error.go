package migrator

import (
	"errors"
	"log"
)

var (
	// Ошибка миграций БД
	ErrMigrator = errors.New("migrator")
	// Пустой провайдер
	ErrNilProvider = errors.New("nil provider")

	// Невалидная конфигурация для сборки
	ErrInvalidBuildConfiguration = errors.New("invalid build configuration")
)

type Logger interface {
	Error(err error)
	Info(s string)
}

func logIfErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func logIfErrFunc(f func() error) {
	logIfErr(f())
}

type simpleLogger struct{}

func (*simpleLogger) Error(err error) {
	if err != nil {
		log.Println(err)
	}
}

func (*simpleLogger) Info(s string) {
	log.Println(s)
}
