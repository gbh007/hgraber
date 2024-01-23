package agentapi

import (
	"app/internal/domain/agent"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (api *API) CreateBooks(ctx context.Context, urls []string) (*agent.CreateBooksResult, error) {
	body, err := json.Marshal(agent.CreateBooksRequest{URLs: urls})
	if err != nil {
		return nil, fmt.Errorf("%s: marshal: %w", apiName, err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, api.url(agent.EndpointBookCreate), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("%s: request: %w", apiName, err)
	}

	data := new(agent.CreateBooksResult)

	err = api.post(ctx, request, data)
	if err != nil {
		return nil, fmt.Errorf("%s: response: %w", apiName, err)
	}

	return data, nil
}
