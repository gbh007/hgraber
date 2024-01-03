package hgraberweb

import (
	"app/internal/controller/hgraberweb/internal/rendering"
	"app/internal/domain/hgraber"
	"net/http"
)

func (ws *WebServer) bookList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			Count int `json:"count"`
			Page  int `json:"page,omitempty"`
		}{}

		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, &request)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)

			return
		}

		data := ws.useCases.GetBooks(ctx, hgraber.BookFilterOuter{
			Count:    request.Count,
			Page:     request.Page,
			NewFirst: true,
		})

		ws.webtool.WriteJSON(ctx, w, http.StatusOK, rendering.BookListResponseFromDomain(ws.outerAddr, data))
	})
}
