package domain

import "errors"

var (
	BookNotFoundError         = errors.New("book not found")
	PageNotFoundError         = errors.New("page not found")
	BookAlreadyExistsError    = errors.New("book already exists")
	UnsupportedAttributeError = errors.New("attribute is not supported")
)
