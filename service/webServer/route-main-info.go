package webServer

import (
	"app/service/webServer/base"
	"net/http"
)

func (ws *WebServer) routeMainInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		stor := ws.Storage

		base.WriteJSON(ctx, w, http.StatusOK, map[string]interface{}{
			"count":               stor.TitlesCount(ctx),
			"not_load_count":      stor.UnloadedTitlesCount(ctx),
			"page_count":          stor.PagesCount(ctx),
			"not_load_page_count": stor.UnloadedPagesCount(ctx),
		})
	})
}
