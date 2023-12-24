package externalfile

import "errors"

var (
	NotFoundError     = errors.New("not found")
	UnauthorizedError = errors.New("unauthorized")
	ForbiddenError    = errors.New("forbidden")
)

const (
	HeaderToken         = "X-Token"
	HeaderBookID        = "X-Book-ID"
	HeaderPageNumber    = "X-Page-Number"
	HeaderPageExtension = "X-Page-Extension"
	HeaderFilename      = "X-Filename"
)

const (
	EndpointPage   = "/page"
	EndpointExport = "/export"
)
