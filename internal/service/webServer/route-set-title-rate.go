package webServer

import (
	"app/pkg/webtool"
	"net/http"
)

func (ws *WebServer) routeSetTitleRate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID   int `json:"id"`
			Rate int `json:"rate"`
		}{}

		ctx := r.Context()

		err := webtool.ParseJSON(r, &request)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		err = ws.storage.UpdateBookRate(ctx, request.ID, request.Rate)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		webtool.WriteJSON(ctx, w, http.StatusOK, struct{}{})
	})
}
