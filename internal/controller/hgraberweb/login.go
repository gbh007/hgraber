package hgraberweb

import (
	"net/http"
)

func (ws *WebServer) login(token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		request := struct {
			Token string `json:"token"`
		}{}

		err := ws.webtool.ParseJSON(r, &request)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)

			return
		}

		if request.Token != token {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, false)

			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     tokenCookieName,
			Value:    token,
			Path:     "/",
			HttpOnly: true,
		})

		ws.webtool.WriteJSON(ctx, w, http.StatusOK, map[string]any{
			"token": token,
		})
	})
}
