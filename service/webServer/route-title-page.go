package webServer

import (
	"app/service/webServer/base"
	"app/service/webServer/rendering"
	"net/http"
)

func (ws *WebServer) routeTitlePage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID   int `json:"id"`
			Page int `json:"page"`
		}{}

		ctx := r.Context()

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		data, err := ws.Storage.GetPage(ctx, request.ID, request.Page)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		base.WriteJSON(ctx, w, http.StatusOK, rendering.PageFullInfoFromStorage(data))
	})
}
