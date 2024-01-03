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
	// Режим отладки
	DebugMode bool
	// Полные пути файлов в логах
	DebugFullpathMode bool
	// Отключить стандартный поток ошибок
	DisableStdErr bool
	// Отключить поток ошибок в файл
	DisableFileErr bool
	// Режим дозаписи файла потока ошибок
	EnableAppendFileErr bool
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

	// потоки логирования
	disableStdErr := flag.Bool("no-stderr", false, "отключить стандартный поток ошибок")
	disableFileErr := flag.Bool("no-stdfile", false, "отключить поток ошибок в файл")
	enableAppendFileErr := flag.Bool("stdfile-append", false, "режим дозаписи файла потока ошибок")

	// размещение данных
	staticDirPath := flag.String("static", "", "папка со статическими файлами")

	// отладка
	debugMode := flag.Bool("debug", false, "активировать режим отладки (дебага)")
	debugFullpathMode := flag.Bool("debug-fullpath", false, "включает длинные пути файлов в логах")

	// агент сервер
	agAddr := flag.String("ag-addr", "", "адрес агент сервера")
	agToken := flag.String("ag-token", "", "токен для доступа к агент серверу")

	flag.Parse()

	return Config{
		Log: LogConfig{
			DebugMode:           *debugMode,
			DebugFullpathMode:   *debugFullpathMode,
			DisableStdErr:       *disableStdErr,
			DisableFileErr:      *disableFileErr,
			EnableAppendFileErr: *enableAppendFileErr,
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
