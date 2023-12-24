package agentapi

import (
	"app/internal/domain/agent"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (api *API) UnprocessedPages(ctx context.Context, limit int) ([]agent.PageToHandle, error) {
	body, err := json.Marshal(agent.UnprocessedRequest{
		Prefixes: api.prefixes,
		Limit:    limit,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: marshal: %w", apiName, err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, api.url(agent.EndpointPageUnprocessed), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("%s: request: %w", apiName, err)
	}

	data := new(agent.UnprocessedResponse[agent.PageToHandle])

	err = api.post(ctx, request, data)
	if err != nil {
		return nil, fmt.Errorf("%s: response: %w", apiName, err)
	}

	return data.ToHandle, nil
}
