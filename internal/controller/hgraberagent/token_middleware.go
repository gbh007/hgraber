package hgraberagent

import (
	"app/internal/domain/agent"
	"net/http"
)

func (c *Controller) tokenMiddleware(next http.Handler) http.Handler {
	// Нет токена, не обрабатываем
	if c.token == "" {
		return next
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userToken := r.Header.Get(agent.HeaderAgentToken)

		if userToken == "" {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		if c.token != userToken {
			w.WriteHeader(http.StatusForbidden)

			return
		}

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
