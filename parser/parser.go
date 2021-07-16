package parser

// Parser интерфейс для реализации парсеров для различных сайтов
type Parser interface {
	Load(URL string) bool
	ParseName() string
	ParsePages() []string
	ParseTags() []string
	ParseAuthors() []string
	ParseCharacters() []string
}

/*


func (p Parser) Load(URL string) bool      { return false}
func (p Parser) ParseName() string         { return "" }
func (p Parser) ParsePages() []string      { return []string{} }
func (p Parser) ParseTags() []string       { return []string{} }
func (p Parser) ParseAuthors() []string    { return []string{} }
func (p Parser) ParseCharacters() []string { return []string{} }

*/
