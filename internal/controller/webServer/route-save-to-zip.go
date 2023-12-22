package webServer

import (
	"app/pkg/webtool"
	"net/http"
)

func (ws *WebServer) routeSaveToZIP() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			From int `json:"from"`
			To   int `json:"to"`
		}{}

		ctx := r.Context()

		err := webtool.ParseJSON(r, &request)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		err = ws.useCases.ExportBooksToZip(ctx, request.From, request.To)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		webtool.WriteJSON(ctx, w, http.StatusOK, struct{}{})
	})
}
