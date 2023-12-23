package webServer

import (
	"net/http"
)

func (ws *WebServer) routeSetTitleRate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID   int `json:"id"`
			Rate int `json:"rate"`
		}{}

		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, &request)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		err = ws.useCases.UpdateBookRate(ctx, request.ID, request.Rate)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		ws.webtool.WriteJSON(ctx, w, http.StatusOK, struct{}{})
	})
}
