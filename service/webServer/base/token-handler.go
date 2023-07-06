package base

import (
	"net/http"
)

const TokenCookieName = "hgraber-access-token"

func TokenHandler(token string, next http.Handler) http.Handler {
	// Нет токена, не обрабатываем
	if token == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(TokenCookieName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		if c.Value != token {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
