package migrator

// migrationFile - файл с миграцией
type migrationFile struct {
	// Номер миграции
	Number int
	// Путь до файла
	Path string
	// Название файла
	Name string
}
