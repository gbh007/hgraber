package webServer

import (
	"app/service/webServer/base"
	"net/http"
)

func (ws *WebServer) routeNewTitle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			URL string `json:"url"`
		}{}

		ctx := r.Context()

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		err = ws.Title.FirstHandle(ctx, request.URL)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusInternalServerError, err)
		} else {
			base.WriteJSON(ctx, w, http.StatusOK, struct{}{})
		}
	})
}