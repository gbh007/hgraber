package agentapi

import (
	"app/internal/domain/agent"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (api *API) SearchBook(ctx context.Context, url string) (*agent.SearchBookIDByURLResponse, error) {
	body, err := json.Marshal(agent.SearchBookIDByURLRequest{URL: url})
	if err != nil {
		return nil, fmt.Errorf("%s: marshal: %w", apiName, err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, api.url(agent.EndpointBookSearch), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("%s: request: %w", apiName, err)
	}

	data := new(agent.SearchBookIDByURLResponse)

	err = api.post(ctx, request, data)
	if err != nil {
		return nil, fmt.Errorf("%s: response: %w", apiName, err)
	}

	return data, nil
}
