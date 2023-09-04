package domain

import "errors"

var (
	TitleNotFoundError        = errors.New("title not found")
	PageNotFoundError         = errors.New("page not found")
	TitleAlreadyExistsError   = errors.New("title already exists")
	UnsupportedAttributeError = errors.New("attribute is not supported")
)
