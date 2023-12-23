package hgraberweb

import (
	"net/http"
)

func (ws *WebServer) routeSaveToZIP() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			From int `json:"from"`
			To   int `json:"to"`
		}{}

		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, &request)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		err = ws.useCases.ExportBooksToZip(ctx, request.From, request.To)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		ws.webtool.WriteJSON(ctx, w, http.StatusOK, struct{}{})
	})
}
