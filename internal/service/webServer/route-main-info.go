package webServer

import (
	"app/internal/service/webServer/base"
	"app/internal/service/webServer/rendering"
	"net/http"
)

func (ws *WebServer) routeMainInfo() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		stor := ws.storage

		base.WriteJSON(ctx, w, http.StatusOK, map[string]interface{}{
			"count":               stor.BooksCount(ctx),
			"not_load_count":      stor.UnloadedBooksCount(ctx),
			"page_count":          stor.PagesCount(ctx),
			"not_load_page_count": stor.UnloadedPagesCount(ctx),
			"monitor":             rendering.MonitorFromDomain(ws.monitor.Info()),
		})
	})
}
