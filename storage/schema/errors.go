package schema

import "errors"

var (
	TitleIndexError     = errors.New("TitleIndexError")
	PageIndexError      = errors.New("PageIndexError")
	TitleDuplicateError = errors.New("TitleDuplicateError")
)
