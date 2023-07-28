package webServer

import (
	"app/config"
	"app/service/webServer/base"
	"app/service/webServer/static"
	"app/super"
	"app/system"
	"context"
	"fmt"
	"net/http"
)

type WebServer struct {
	Storage   super.Storage
	Title     super.TitleHandler
	Page      super.PageHandler
	Addr      string
	StaticDir string
	Token     string
}

func Init(
	storage super.Storage,
	title super.TitleHandler,
	page super.PageHandler,
	config config.WebServerConfig,
) *WebServer {
	return &WebServer{
		Storage:   storage,
		Title:     title,
		Page:      page,
		Addr:      fmt.Sprintf("%s:%d", config.Host, config.Port),
		StaticDir: config.StaticDirPath,
		Token:     config.Token,
	}
}

func makeServer(parentCtx context.Context, ws *WebServer) *http.Server {
	mux := http.NewServeMux()

	// обработчик статики
	if ws.StaticDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(ws.StaticDir)))
	} else {
		mux.Handle("/", http.FileServer(http.FS(static.StaticDir)))
	}

	// обработчик файлов
	mux.Handle("/file/", base.TokenHandler(ws.Token,
		http.StripPrefix(
			"/file/",
			http.FileServer(http.Dir(system.GetFileStoragePath(parentCtx))),
		),
	))

	// API
	mux.Handle("/auth/login", ws.routeLogin(ws.Token))
	mux.Handle("/info", base.TokenHandler(ws.Token, ws.routeMainInfo()))
	mux.Handle("/new", base.TokenHandler(ws.Token, ws.routeNewTitle()))
	mux.Handle("/title/list", base.TokenHandler(ws.Token, ws.routeTitleList()))
	mux.Handle("/title/details", base.TokenHandler(ws.Token, ws.routeTitleInfo()))
	mux.Handle("/title/page", base.TokenHandler(ws.Token, ws.routeTitlePage()))
	mux.Handle("/to-zip", base.TokenHandler(ws.Token, ws.routeSaveToZIP()))
	mux.Handle("/app/info", base.TokenHandler(ws.Token, ws.routeAppInfo()))
	mux.Handle("/title/rate", base.TokenHandler(ws.Token, ws.routeSetTitleRate()))
	mux.Handle("/title/page/rate", base.TokenHandler(ws.Token, ws.routeSetPageRate()))

	server := &http.Server{
		Addr: ws.Addr,
		Handler: base.PanicDefender(
			base.Stopwatch(mux),
		),
		ErrorLog:    system.StdErrorLogger(parentCtx),
		BaseContext: base.NewBaseContext(parentCtx),
	}

	return server
}
