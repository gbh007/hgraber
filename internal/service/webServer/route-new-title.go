package webServer

import (
	"app/internal/service/webServer/base"
	"app/internal/service/webServer/rendering"
	"net/http"
)

func (ws *WebServer) routeNewTitle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := struct {
			URL  string   `json:"url,omitempty"`
			URLs []string `json:"urls,omitempty"`
		}{}

		ctx := r.Context()

		err := base.ParseJSON(r, &request)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusBadRequest, err)
			return
		}

		if len(request.URLs) > 0 {
			base.WriteJSON(
				ctx, w, http.StatusOK,
				rendering.HandleMultipleResultFromDomain(ws.title.FirstHandleMultiple(ctx, request.URLs)),
			)

			return
		}

		err = ws.title.FirstHandle(ctx, request.URL)
		if err != nil {
			base.WriteJSON(ctx, w, http.StatusInternalServerError, err)
		} else {
			base.WriteJSON(ctx, w, http.StatusOK, struct{}{})
		}
	})
}
