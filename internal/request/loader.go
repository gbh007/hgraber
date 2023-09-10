package request

import (
	"app/system"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"
)

const defaultUserAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"

type Requester struct {
	client *http.Client
}

func New() *Requester {
	return &Requester{
		client: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// requestBuffer запрашивает данные по урле и возвращает их в виде буффера
func (r *Requester) requestBuffer(ctx context.Context, URL string) (bytes.Buffer, error) {
	buff := bytes.Buffer{}
	req, err := http.NewRequest(http.MethodGet, URL, &bytes.Buffer{})
	if err != nil {
		system.Error(ctx, err)
		return buff, err
	}

	req.Header.Set("User-Agent", defaultUserAgent)

	// выполняем запрос
	response, err := r.client.Do(req)

	if err != nil {
		system.Error(ctx, err)
		return buff, err
	}

	defer system.IfErrFunc(ctx, response.Body.Close)

	if response.StatusCode < 200 || response.StatusCode > 299 {
		err = fmt.Errorf("%s ошибка %s", URL, response.Status)
		system.Error(ctx, err)

		return buff, err
	}

	_, err = buff.ReadFrom(response.Body)
	if err != nil {
		system.Error(ctx, err)

		return buff, err
	}

	return buff, nil
}

// RequestString запрашивает данные по урле и возвращает их строкой
func (r *Requester) RequestString(ctx context.Context, URL string) (string, error) {
	buff, err := r.requestBuffer(ctx, URL)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}

// RequestBytes запрашивает данные по урле и возвращает их массивом байт
func (r *Requester) RequestBytes(ctx context.Context, URL string) ([]byte, error) {
	buff, err := r.requestBuffer(ctx, URL)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
