package base

import (
	"app/system"
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
