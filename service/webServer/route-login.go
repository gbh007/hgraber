package webServer

import (
	"app/service/webServer/base"
	"net/http"
)

func (ws *WebServer) routeLogin(token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			Token string `json:"token"`
		}{}
		ctx := r.Context()

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusBadRequest, err)

			return
		}

		if request.Token != token {
			base.WriteJSON(ctx, w, http.StatusBadRequest, false)

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     base.TokenCookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		})

		base.WriteJSON(ctx, w, http.StatusOK, struct{}{})
	})
}
