package webServer

import (
	"app/internal/service/webServer/base"
	"app/internal/service/webServer/rendering"
	"net/http"
)

func (ws *WebServer) routeTitleInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			ID int `json:"id"`
		}{}
		ctx := r.Context()

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		data, err := ws.storage.GetBook(ctx, request.ID)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusInternalServerError, err)
			return
		}

		base.WriteJSON(ctx, w, http.StatusOK, rendering.TitleFromStorage(data))
	})
}
