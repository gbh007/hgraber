package domain

type Attribute string

var AllAttributes = []Attribute{
	AttrAuthor,
	AttrCategory,
	AttrCharacter,
	AttrGroup,
	AttrLanguage,
	AttrParody,
	AttrTag,
}

const (
	AttrAuthor    Attribute = "author"
	AttrCategory  Attribute = "category"
	AttrCharacter Attribute = "character"
	AttrGroup     Attribute = "group"
	AttrLanguage  Attribute = "language"
	AttrParody    Attribute = "parody"
	AttrTag       Attribute = "tag"
)
