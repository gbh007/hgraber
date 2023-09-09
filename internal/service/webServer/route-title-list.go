package webServer

import (
	"app/internal/domain"
	"app/internal/service/webServer/base"
	"app/internal/service/webServer/rendering"
	"net/http"
)

func (ws *WebServer) routeTitleList() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			Count  int `json:"count"`
			Offset int `json:"offset,omitempty"`
		}{}

		ctx := r.Context()

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		data := ws.storage.GetTitles(ctx, domain.BookFilter{
			Limit:    request.Count,
			Offset:   request.Offset,
			NewFirst: true,
		})

		base.WriteJSON(ctx, w, http.StatusOK, rendering.TitlesFromStorage(data))
	})
}
