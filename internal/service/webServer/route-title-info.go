package webServer

import (
	"app/internal/service/webServer/rendering"
	"app/pkg/webtool"
	"net/http"
)

func (ws *WebServer) routeTitleInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID int `json:"id"`
		}{}
		ctx := r.Context()

		err := webtool.ParseJSON(r, &request)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		data, err := ws.storage.GetBook(ctx, request.ID)
		if err != nil {
			webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		webtool.WriteJSON(ctx, w, http.StatusOK, rendering.TitleFromStorageWrap(ws.outerAddr)(data))
	})
}
