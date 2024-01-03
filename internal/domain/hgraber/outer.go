package hgraber

// FIXME: требуется перенос структур из этого файла, предварительно в новый домен.
// Альтернативно можно часть структур из этого домена перенести в новый.

// FIXME: слить с другим фильтром, добавить методы limit/offset как бизнес логику
type BookFilterOuter struct {
	Count    int
	Page     int
	NewFirst bool
}

type FilteredBooks struct {
	Books       []Book
	Pages       []int
	CurrentPage int
}
