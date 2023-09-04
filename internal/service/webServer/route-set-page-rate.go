package webServer

import (
	"app/internal/service/webServer/base"
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

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		err = ws.storage.UpdatePageRate(ctx, request.ID, request.Page, request.Rate)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		base.WriteJSON(ctx, w, http.StatusOK, struct{}{})
	})
}
