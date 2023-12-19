package webServer

import (
	"app/internal/domain"
	"app/internal/service/webServer/rendering"
	"app/pkg/webtool"
	"net/http"
)

func (ws *WebServer) routeTitleList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			Count  int `json:"count"`
			Offset int `json:"offset,omitempty"`
		}{}

		ctx := r.Context()

		err := webtool.ParseJSON(r, &request)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		data := ws.storage.GetBooks(ctx, domain.BookFilter{
			Limit:    request.Count,
			Offset:   request.Offset,
			NewFirst: true,
		})

		webtool.WriteJSON(ctx, w, http.StatusOK, rendering.TitlesFromStorage(ws.outerAddr, data))
	})
}
