package webtool

import (
	"app/system"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
)

func ParseJSON(r *http.Request, data interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		system.Debug(r.Context(), err)
	}

	return err
}

func NewBaseContext(ctx context.Context) func(l net.Listener) context.Context {
	return func(l net.Listener) context.Context { return system.NewUserContext(ctx) }
}

func WriteJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	enc := json.NewEncoder(w)

	if errData, ok := data.(error); ok {
		data = errData.Error()
	}

	if system.IsDebug(ctx) {
		enc.SetIndent("", "  ")
	}

	if err := enc.Encode(data); err != nil {
		system.Error(ctx, err)
	}

}


func WriteNoContent(ctx context.Context, w http.ResponseWriter, ) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.WriteHeader(http.StatusNoContent)
}

func WritePlain(ctx context.Context, w http.ResponseWriter, statusCode int, data string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)

	_, err := io.WriteString(w, data)
	system.IfErr(ctx, err)
}
