package parser

import (
	"app/system"
	"bytes"
	"fmt"
	"net/http"
	"time"
)

// requestBuffer запрашивает данные по урле и возвращает их в виде буффера
func requestBuffer(ctx system.Context, URL string) (bytes.Buffer, error) {
	buff := bytes.Buffer{}
	req, err := http.NewRequest(http.MethodGet, URL, &bytes.Buffer{})
	if err != nil {
		system.Error(ctx, err)
		return buff, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36")
	req.Close = true
	// выполняем запрос
	response, err := (&http.Client{
		Timeout: time.Minute,
		// Transport: &http.Transport{
		// 	TLSClientConfig: &tls.Config{
		// 		InsecureSkipVerify: true,
		// 	},
		// },
	}).Do(req)
	if err != nil {
		system.Error(ctx, err)
		return buff, err
	}
	defer response.Body.Close()
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
func RequestString(ctx system.Context, URL string) (string, error) {
	buff, err := requestBuffer(ctx, URL)
	return buff.String(), err
}

// RequestBytes запрашивает данные по урле и возвращает их массивом байт
func RequestBytes(ctx system.Context, URL string) ([]byte, error) {
	buff, err := requestBuffer(ctx, URL)
	return buff.Bytes(), err
}
