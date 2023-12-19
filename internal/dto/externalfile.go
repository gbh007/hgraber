package dto

import "errors"

var (
	ExternalFileNotFoundError     = errors.New("not found")
	ExternalFileUnauthorizedError = errors.New("unauthorized")
	ExternalFileForbiddenError    = errors.New("forbidden")
)

const (
	ExternalFileToken         = "X-Token"
	ExternalFileBookID        = "X-Book-ID"
	ExternalFilePageNumber    = "X-Page-Number"
	ExternalFilePageExtension = "X-Page-Extension"
	ExternalFileFilename      = "X-Filename"
)

const (
	ExternalFileEndpointPage   = "/page"
	ExternalFileEndpointExport = "/export"
)
