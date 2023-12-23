package webServer

import (
	"app/internal/controller/webServer/internal/rendering"
	"net/http"
)

func (ws *WebServer) routeTitleInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID int `json:"id"`
		}{}
		ctx := r.Context()

		err := ws.webtool.ParseJSON(r, &request)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		data, err := ws.useCases.GetBook(ctx, request.ID)
		if err != nil {
			ws.webtool.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		ws.webtool.WriteJSON(ctx, w, http.StatusOK, rendering.TitleFromStorageWrap(ws.outerAddr)(data))
	})
}
