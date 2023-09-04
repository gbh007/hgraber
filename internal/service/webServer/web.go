package webServer

import (
	"app/internal/config"
	"app/internal/domain"
	"app/internal/service/webServer/base"
	"app/internal/service/webServer/static"
	"app/system"
	"context"
	"fmt"
	"net/http"
)

type pageHandler interface {
	ExportTitlesToZip(ctx context.Context, from, to int) error
}

type titleHandler interface {
	// FirstHandle обрабатывает данные тайтла (новое добавление, упрощенное без парса страниц)
	FirstHandle(ctx context.Context, u string) error
}

type storage interface {
	GetPage(ctx context.Context, id int, page int) (*domain.PageFullInfo, error)
	GetTitle(ctx context.Context, id int) (domain.Title, error)
	GetTitles(ctx context.Context, offset int, limit int) []domain.Title
	PagesCount(ctx context.Context) int
	TitlesCount(ctx context.Context) int
	UnloadedPagesCount(ctx context.Context) int
	UnloadedTitlesCount(ctx context.Context) int
	UpdatePageRate(ctx context.Context, id int, page int, rate int) error
	UpdateTitleRate(ctx context.Context, id int, rate int) error
}

type WebServer struct {
	storage   storage
	title     titleHandler
	page      pageHandler
	addr      string
	staticDir string
	token     string
}

func Init(
	storage storage,
	title titleHandler,
	page pageHandler,
	config config.WebServerConfig,
) *WebServer {
	return &WebServer{
		storage:   storage,
		title:     title,
		page:      page,
		addr:      fmt.Sprintf("%s:%d", config.Host, config.Port),
		staticDir: config.StaticDirPath,
		token:     config.Token,
	}
}

func makeServer(parentCtx context.Context, ws *WebServer) *http.Server {
	mux := http.NewServeMux()

	// обработчик статики
	if ws.staticDir != "" {
		mux.Handle("/", http.FileServer(http.Dir(ws.staticDir)))
	} else {
		mux.Handle("/", http.FileServer(http.FS(static.StaticDir)))
	}

	// обработчик файлов
	mux.Handle("/file/", base.TokenHandler(ws.token,
		http.StripPrefix(
			"/file/",
			http.FileServer(http.Dir(system.GetFileStoragePath(parentCtx))),
		),
	))

	// API
	mux.Handle("/auth/login", ws.routeLogin(ws.token))
	mux.Handle("/info", base.TokenHandler(ws.token, ws.routeMainInfo()))
	mux.Handle("/new", base.TokenHandler(ws.token, ws.routeNewTitle()))
	mux.Handle("/title/list", base.TokenHandler(ws.token, ws.routeTitleList()))
	mux.Handle("/title/details", base.TokenHandler(ws.token, ws.routeTitleInfo()))
	mux.Handle("/title/page", base.TokenHandler(ws.token, ws.routeTitlePage()))
	mux.Handle("/to-zip", base.TokenHandler(ws.token, ws.routeSaveToZIP()))
	mux.Handle("/app/info", base.TokenHandler(ws.token, ws.routeAppInfo()))
	mux.Handle("/title/rate", base.TokenHandler(ws.token, ws.routeSetTitleRate()))
	mux.Handle("/title/page/rate", base.TokenHandler(ws.token, ws.routeSetPageRate()))

	server := &http.Server{
		Addr: ws.addr,
		Handler: base.PanicDefender(
			base.Stopwatch(mux),
		),
		ErrorLog:    system.StdErrorLogger(parentCtx),
		BaseContext: base.NewBaseContext(parentCtx),
	}

	return server
}
