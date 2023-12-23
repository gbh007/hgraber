package hgraberweb

import (
	"app/internal/controller/hgraberweb/internal/rendering"
	"net/http"
)

func (ws *WebServer) routeTitlePage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID   int `json:"id"`
			Page int `json:"page"`
		}{}

		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, &request)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		data, err := ws.useCases.GetPage(ctx, request.ID, request.Page)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		ws.webtool.WriteJSON(ctx, w, http.StatusOK, rendering.PageFromStorageWrap(ws.outerAddr)(*data))
	})
}
