package webServer

import (
	"app/pkg/webtool"
	"net/http"
)

func (ws *WebServer) routeSetPageRate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID   int `json:"id"`
			Page int `json:"page"`
			Rate int `json:"rate"`
		}{}

		ctx := r.Context()

		err := webtool.ParseJSON(r, &request)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		err = ws.storage.UpdatePageRate(ctx, request.ID, request.Page, request.Rate)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		webtool.WriteJSON(ctx, w, http.StatusOK, struct{}{})
	})
}
