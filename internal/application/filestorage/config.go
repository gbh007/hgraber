package filestorage

import "flag"

type config struct {
	LoadPath   string
	ExportPath string

	ReadOnly bool

	Addr  string
	Token string

	Debug bool
	Trace bool
}

func parseFlag() config {
	// веб сервер
	addr := flag.String("addr", ":8080", "адрес веб сервера")
	token := flag.String("token", "", "токен для доступа к контенту")

	// размещение данных
	fileStoragePath := flag.String("fs", "loads", "директория для данных")
	fileExportPath := flag.String("fe", "exported", "директория для экспорта файлов")

	readOnly := flag.Bool("read-only", false, "режим только просмотра")

	// Отладка
	debug := flag.Bool("debug", false, "Режим отладки")
	debugTrace := flag.Bool("debug-trace", false, "Режим стектрейсов")

	flag.Parse()

	return config{
		LoadPath:   *fileStoragePath,
		ExportPath: *fileExportPath,

		ReadOnly: *readOnly,

		Addr:  *addr,
		Token: *token,

		Debug: *debug,
		Trace: *debugTrace,
	}
}
