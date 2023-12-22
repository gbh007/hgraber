package webtool

import (
	"errors"
	"net/http"
)

var errPanicDetected = errors.New("нарушение потока выполнения запроса")

func PanicDefender(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			p := recover()
			if p != nil {
				// FIXME: поддержать
				// system.Warning(r.Context(), "обнаружена паника", p)

				WriteJSON(r.Context(), w, http.StatusInternalServerError, errPanicDetected)
			}
		}()
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(w, r)
	})
}
