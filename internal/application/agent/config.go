package agent

import "flag"

type configRaw struct {
	Scheme string
	Addr   string
	Token  string
	Name   string

	Debug bool
	Trace bool
}

func parseFlag() configRaw {
	hgScheme := flag.String("scheme", "http", "Схема соединения с hgraber")
	hgAddr := flag.String("addr", "localhost:8080", "Адрес соединения с hgraber")
	hgToken := flag.String("token", "", "Токен для доступа к соединению с hgraber")
	name := flag.String("name", "simple-agent", "Название агента")

	// Отладка
	debug := flag.Bool("debug", false, "Режим отладки")
	debugTrace := flag.Bool("debug-trace", false, "Режим стектрейсов")

	flag.Parse()

	cfg := configRaw{
		Scheme: *hgScheme,
		Addr:   *hgAddr,
		Token:  *hgToken,
		Name:   *name,

		Debug: *debug,
		Trace: *debugTrace,
	}

	return cfg
}
