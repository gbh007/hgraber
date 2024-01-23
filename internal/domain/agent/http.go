package agent

import "errors"

var (
	UnauthorizedError = errors.New("unauthorized")
	ForbiddenError    = errors.New("forbidden")
)

const (
	HeaderAgentToken    = "X-Agent-Token"
	HeaderAgentName     = "X-Agent-Name"
	HeaderBookID        = "X-Book-ID"
	HeaderPageNumber    = "X-Page-Number"
	HeaderPageExtension = "X-Page-Extension"
	HeaderPageUrl       = "X-Page-Url"
)

const (
	EndpointBookUnprocessed = "/book/unprocessed"
	EndpointBookUpdate      = "/book/update"
	EndpointBookCreate      = "/book/create"

	EndpointPageUnprocessed = "/page/unprocessed"
	EndpointPageUpload      = "/page/upload"
)

type UnprocessedRequest struct {
	Prefixes []string `json:"prefixes"`
	Limit    int      `json:"limit"`
}

type UnprocessedResponse[T BookToHandle | PageToHandle] struct {
	ToHandle []T `json:"to_handle,omitempty"`
}

type CreateBooksRequest struct {
	URLs []string `json:"urls"`
}
