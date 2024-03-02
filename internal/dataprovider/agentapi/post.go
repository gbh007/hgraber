package agentapi

import (
	"app/internal/domain/agent"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (api *API) post(ctx context.Context, request *http.Request, dataToUnmarshal any) error {
	request.Header.Set(agent.HeaderAgentToken, api.token)
	request.Header.Set(agent.HeaderAgentName, api.agentName)

	response, err := api.client.Do(request)
	if err != nil {
		return fmt.Errorf("%s: do: %w", apiName, err)
	}

	defer func() {
		bodyCloseErr := response.Body.Close()
		if bodyCloseErr != nil {
			api.logger.ErrorContext(ctx, bodyCloseErr.Error())
		}
	}()

	switch response.StatusCode {
	case http.StatusNoContent:
		return nil

	case http.StatusOK:
		if dataToUnmarshal == nil {
			return nil
		}

		err = json.NewDecoder(response.Body).Decode(dataToUnmarshal)
		if err != nil {
			return fmt.Errorf("%s: decode: %w", apiName, err)
		}

		return nil

	case http.StatusUnauthorized:
		return fmt.Errorf("%s: %w", apiName, agent.UnauthorizedError)
	case http.StatusForbidden:
		return fmt.Errorf("%s: %w", apiName, agent.ForbiddenError)
	}

	partOfBodyData := make([]byte, 100)
	n, err := response.Body.Read(partOfBodyData)
	if err != nil {
		return fmt.Errorf("%s: read code(%d) : %w", apiName, response.StatusCode, err)
	}

	if n < len(partOfBodyData) {
		partOfBodyData = partOfBodyData[:n]
	}

	return fmt.Errorf("%s: unknown error code(%d): %s", apiName, response.StatusCode, string(partOfBodyData))
}

func (api *API) url(path string) string {
	u := url.URL{
		Scheme: api.scheme,
		Host:   api.hostWithPort,
		Path:   path,
	}

	return u.String()
}
