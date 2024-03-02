package inmemory

import "flag"

// Config - конфигурация приложения.
type Config struct {
	// Конфигурация веб сервера
	WebServer WebServerConfig
	// Конфигурация логов приложения
	Log LogConfig
	// Агент сервер
	Ag configAg
}

// WebServerConfig - конфигурация веб сервера.
type WebServerConfig struct {
	// Хост веб сервера
	Host string
	// Порт веб сервера
	Port int
	// Токен авторизации веб сервера
	Token string
	// Путь до папки со статическими файлами.
	StaticDirPath string
}

// LogConfig - конфигурация логов приложения.
type LogConfig struct {
	Debug bool
	Trace bool
}

type configAg struct {
	Addr  string
	Token string
}

func parseFlag() Config {
	// базовые опции
	webPort := flag.Int("p", 8080, "порт веб сервера")
	webHost := flag.String("h", "", "хост веб сервера")
	token := flag.String("access-token", "", "токен для доступа к контенту")

	// размещение данных
	staticDirPath := flag.String("static", "", "папка со статическими файлами")

	// Отладка
	debug := flag.Bool("debug", false, "Режим отладки")
	debugTrace := flag.Bool("debug-trace", false, "Режим стектрейсов")

	// агент сервер
	agAddr := flag.String("ag-addr", "", "адрес агент сервера")
	agToken := flag.String("ag-token", "", "токен для доступа к агент серверу")

	flag.Parse()

	return Config{
		Log: LogConfig{
			Debug: *debug,
			Trace: *debugTrace,
		},
		WebServer: WebServerConfig{
			Host:          *webHost,
			Port:          *webPort,
			Token:         *token,
			StaticDirPath: *staticDirPath,
		},
		Ag: configAg{
			Addr:  *agAddr,
			Token: *agToken,
		},
	}
}
