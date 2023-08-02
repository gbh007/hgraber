package schema

import "errors"

var (
	TitleNotFoundError      = errors.New("title not found")
	PageNotFoundError       = errors.New("page not found")
	TitleAlreadyExistsError = errors.New("title already exists")
)
