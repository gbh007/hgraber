package simple

import "flag"

// Config - конфигурация приложения.
type Config struct {
	// Базовая конфигурация приложения
	Base BaseConfig
	// Конфигурация веб сервера
	WebServer WebServerConfig
	// Конфигурация логов приложения
	Log LogConfig
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

// BaseConfig - базовая конфигурация приложения.
type BaseConfig struct {
	// Режим просмотра
	OnlyView bool
	// Путь до каталога с файлами (изображениями)
	FileStoragePath string
	// Путь для каталога для экспорта файлов
	FileExportPath string
	// Путь до файла базы
	DBFilePath string
}

// LogConfig - конфигурация логов приложения.
type LogConfig struct {
	Debug bool
	Trace bool
}

func parseFlag() Config {
	// базовые опции
	webPort := flag.Int("p", 8080, "порт веб сервера")
	webHost := flag.String("h", "", "хост веб сервера")
	onlyView := flag.Bool("v", false, "режим только просмотра")
	token := flag.String("access-token", "", "токен для доступа к контенту")

	// размещение данных
	fileStoragePath := flag.String("fs", "loads", "директория для данных")
	fileExportPath := flag.String("fe", "exported", "директория для экспорта файлов")
	dbFilePath := flag.String("db", "db.json", "файл базы")
	staticDirPath := flag.String("static", "", "папка со статическими файлами")

	// Отладка
	debug := flag.Bool("debug", false, "Режим отладки")
	debugTrace := flag.Bool("debug-trace", false, "Режим стектрейсов")

	flag.Parse()

	return Config{
		Base: BaseConfig{
			OnlyView:        *onlyView,
			FileStoragePath: *fileStoragePath,
			FileExportPath:  *fileExportPath,
			DBFilePath:      *dbFilePath,
		},
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
	}
}
