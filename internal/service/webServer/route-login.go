package webServer

import (
	"app/pkg/webtool"
	"net/http"
)

func (ws *WebServer) routeLogin(token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			Token string `json:"token"`
		}{}
		ctx := r.Context()

		err := webtool.ParseJSON(r, &request)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)

			return
		}

		if request.Token != token {
			webtool.WriteJSON(ctx, w, http.StatusBadRequest, false)

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     tokenCookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		})

		webtool.WriteJSON(ctx, w, http.StatusOK, struct{}{})
	})
}
