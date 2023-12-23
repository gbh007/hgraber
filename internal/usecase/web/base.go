package web

import (
	"app/pkg/ctxtool"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
)

func (uc *UseCase) ParseJSON(r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		uc.logger.Debug(r.Context(), err)
	}

	return err
}

func (uc *UseCase) NewBaseContext(ctx context.Context) func(l net.Listener) context.Context {
	return func(l net.Listener) context.Context { return ctxtool.NewUserContext(ctx) }
}

func (uc *UseCase) WriteJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)

	if errData, ok := data.(error); ok {
		data = errData.Error()
	}

	if uc.debug {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(data); err != nil {
		uc.logger.Error(ctx, err)
	}

}

func (uc *UseCase) WriteNoContent(ctx context.Context, w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (uc *UseCase) WritePlain(ctx context.Context, w http.ResponseWriter, statusCode int, data string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)

	_, err := io.WriteString(w, data)
	uc.logger.IfErr(ctx, err)
}
