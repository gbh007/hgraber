package base

import (
	"app/system"
	"fmt"
	"net/http"
)

func PanicDefender(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			p := recover()
			if p != nil {
				system.Warning(r.Context(), "обнаружена паника", p)

				WriteJSON(r.Context(), w, http.StatusInternalServerError, ErrPanicDetected)
			}
		}()
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func Stopwatch(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer system.Stopwatch(r.Context(), fmt.Sprintf("ws path %s", r.URL.Path))()

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
