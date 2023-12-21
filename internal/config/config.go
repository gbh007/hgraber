package config

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
	// Тип БД
	DBType string
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

func ParseFlag() Config {
	// базовые опции
	webPort := flag.Int("p", 8080, "порт веб сервера")
	webHost := flag.String("h", "", "хост веб сервера")
	onlyView := flag.Bool("v", false, "режим только просмотра")
	token := flag.String("access-token", "", "токен для доступа к контенту")

	// потоки логирования
	disableStdErr := flag.Bool("no-stderr", false, "отключить стандартный поток ошибок")
	disableFileErr := flag.Bool("no-stdfile", false, "отключить поток ошибок в файл")
	enableAppendFileErr := flag.Bool("stdfile-append", false, "режим дозаписи файла потока ошибок")

	// размещение данных
	fileStoragePath := flag.String("fs", "loads", "директория для данных")
	fileExportPath := flag.String("fe", "exported", "директория для экспорта файлов")
	dbFilePath := flag.String("db", "db.json", "файл базы")
	dbType := flag.String("db-type", "jdb", "Тип БД: jdb, pg")
	staticDirPath := flag.String("static", "", "папка со статическими файлами")

	// отладка
	debugMode := flag.Bool("debug", false, "активировать режим отладки (дебага)")
	debugFullpathMode := flag.Bool("debug-fullpath", false, "включает длинные пути файлов в логах")

	flag.Parse()

	return Config{
		Base: BaseConfig{
			OnlyView:        *onlyView,
			FileStoragePath: *fileStoragePath,
			FileExportPath:  *fileExportPath,
			DBFilePath:      *dbFilePath,
			DBType:          *dbType,
		},
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
	}
}
