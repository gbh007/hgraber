package webServer

import "net/http"

const tokenCookieName = "hgraber-access-token"

func tokenHandler(token string, next http.Handler) http.Handler {
	// Нет токена, не обрабатываем
	if token == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
