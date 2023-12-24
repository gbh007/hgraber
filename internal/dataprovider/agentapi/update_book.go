package agentapi

import (
	"app/internal/domain/agent"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (api *API) UpdateBook(ctx context.Context, book agent.BookToUpdate) error {
	body, err := json.Marshal(book)
	if err != nil {
		return fmt.Errorf("%s: marshal: %w", apiName, err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, api.url(agent.EndpointBookUpdate), bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("%s: request: %w", apiName, err)
	}

	err = api.post(ctx, request, nil)
	if err != nil {
		return fmt.Errorf("%s: response: %w", apiName, err)
	}

	return nil
}
