package webServer

import (
	"app/internal/controller/webServer/internal/rendering"
	"app/internal/domain"
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

		data := ws.useCases.GetBooks(ctx, domain.BookFilter{
			Limit:    request.Count,
			Offset:   request.Offset,
			NewFirst: true,
		})

		webtool.WriteJSON(ctx, w, http.StatusOK, rendering.TitlesFromStorage(ws.outerAddr, data))
	})
}
