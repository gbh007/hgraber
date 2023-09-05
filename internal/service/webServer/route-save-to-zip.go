package webServer

import (
	"app/internal/service/webServer/base"
	"net/http"
)

func (ws *WebServer) routeSaveToZIP() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			From int `json:"from"`
			To   int `json:"to"`
		}{}

		ctx := r.Context()

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		err = ws.page.ExportTitlesToZip(ctx, request.From, request.To)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		base.WriteJSON(ctx, w, http.StatusOK, struct{}{})
	})
}