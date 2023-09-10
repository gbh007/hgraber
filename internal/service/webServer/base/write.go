package base

import (
	"app/system"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

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

func WritePlain(ctx context.Context, w http.ResponseWriter, statusCode int, data string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(statusCode)

	_, err := io.WriteString(w, data)
	system.IfErr(ctx, err)
}
