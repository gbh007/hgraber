package server

import (
	"flag"
	"os"
)

type configRaw struct {
	fs configFS
	ws configWS
	ag configAg

	ReadOnly bool

	PGSource string
}

type configFS struct {
	Scheme string
	Addr   string
	Token  string
}

type configAg struct {
	Addr  string
	Token string
}

type configWS struct {
	Addr      string
	OuterAddr string
	Token     string
	Static    string
}

func parseFlag() configRaw {
	// файловое хранилище
	fsScheme := flag.String("fs-scheme", "http", "Схема соединения с файловой системой")
	fsAddr := flag.String("fs-addr", "localhost:8080", "Адрес соединения с файловой системой")
	fsToken := flag.String("fs-token", "", "Токен для доступа к ресурсам соединения с файловой системой")

	pgSource := flag.String("pg-source", "", "Строка подключения к PostgreSQL")
	onlyView := flag.Bool("v", false, "режим только просмотра")

	// веб сервер
	wsAddr := flag.String("ws-addr", ":8080", "адрес веб сервера")
	wsOuterAddr := flag.String("ws-outer-addr", "http://localhost:8080", "Внешний адрес веб сервера")
	wsToken := flag.String("ws-token", "", "токен для доступа к контенту")
	wsStatic := flag.String("ws-static", "", "папка со статическими файлами")

	// агент сервер
	agAddr := flag.String("ag-addr", "", "адрес агент сервера")
	agToken := flag.String("ag-token", "", "токен для доступа к агент серверу")

	flag.Parse()

	cfg := configRaw{
		fs: configFS{
			Scheme: *fsScheme,
			Addr:   *fsAddr,
			Token:  *fsToken,
		},
		ws: configWS{
			Addr:      *wsAddr,
			Token:     *wsToken,
			Static:    *wsStatic,
			OuterAddr: *wsOuterAddr,
		},
		ag: configAg{
			Addr:  *agAddr,
			Token: *agToken,
		},
		PGSource: *pgSource,
		ReadOnly: *onlyView,
	}

	if pgSource := os.Getenv("PG_SOURCE"); pgSource != "" {
		cfg.PGSource = pgSource
	}

	return cfg
}
