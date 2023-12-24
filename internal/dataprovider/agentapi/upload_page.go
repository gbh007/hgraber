package agentapi

import (
	"app/internal/domain/agent"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func (api *API) UploadPage(ctx context.Context, info agent.PageInfoToUpload, body io.Reader) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, api.url(agent.EndpointPageUpload), body)
	if err != nil {
		return fmt.Errorf("%s: request: %w", apiName, err)
	}

	request.Header.Set(agent.HeaderBookID, strconv.Itoa(info.BookID))
	request.Header.Set(agent.HeaderPageNumber, strconv.Itoa(info.PageNumber))
	request.Header.Set(agent.HeaderPageExtension, info.Ext)
	request.Header.Set(agent.HeaderPageUrl, info.URL)

	return api.post(ctx, request, nil)
}
