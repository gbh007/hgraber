package hgraberweb

import (
	"net/http"
)

const (
	tokenCookieName = "hgraber-access-token"
	tokenHeaderName = "X-Token"
)

func tokenHandler(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler { // Спешиал фор Игорь
		// Нет токена, не обрабатываем
		if token == "" {
			return next
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userToken := r.Header.Get(tokenHeaderName)

			if userToken != "" && token == userToken {
				next.ServeHTTP(w, r)

				return
			}

			c, err := r.Cookie(tokenCookieName)
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
}
