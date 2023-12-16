package webServer

import (
	"app/internal/config"
	"app/internal/domain"
	"app/internal/service/webServer/base"
	"app/internal/service/webServer/static"
	"app/system"
	"context"
	"fmt"
	"io"
	"net/http"
)

type pageHandler interface {
	ExportBooksToZip(ctx context.Context, from, to int) error
}

type titleHandler interface {
	// FirstHandle обрабатывает данные тайтла (новое добавление, упрощенное без парса страниц)
	FirstHandle(ctx context.Context, u string) error
}

type storage interface {
	GetPage(ctx context.Context, id int, page int) (*domain.Page, error)
	GetBook(ctx context.Context, id int) (domain.Book, error)
	GetBooks(ctx context.Context, filter domain.BookFilter) []domain.Book
	PagesCount(ctx context.Context) int
	BooksCount(ctx context.Context) int
	UnloadedPagesCount(ctx context.Context) int
	UnloadedBooksCount(ctx context.Context) int
	UpdatePageRate(ctx context.Context, id int, page int, rate int) error
	UpdateBookRate(ctx context.Context, id int, rate int) error
}

type files interface {
	OpenPageFile(ctx context.Context, id, page int, ext string) (io.ReadCloser, error)
}

type monitor interface {
	Info() []domain.MonitorStat
}

type WebServer struct {
	storage storage
	title   titleHandler
	page    pageHandler
	files   files
	monitor monitor

	addr      string
	outerAddr string
	staticDir string
	token     string
}

func Init(
	storage storage,
	title titleHandler,
	page pageHandler,
	files files,
	monitor monitor,
	config config.WebServerConfig,
) *WebServer {
	return &WebServer{
		storage: storage,
		title:   title,
		page:    page,
		files:   files,
		monitor: monitor,

		addr:      fmt.Sprintf("%s:%d", config.Host, config.Port),
		outerAddr: fmt.Sprintf("http://%s:%d", config.Host, config.Port),
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
	mux.Handle("/file/", base.TokenHandler(ws.token, http.StripPrefix("/file/", ws.getFile())))

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
			base.Stopwatch(
				base.CORS(mux),
			),
		),
		ErrorLog:    system.StdErrorLogger(parentCtx),
		BaseContext: base.NewBaseContext(context.WithoutCancel(parentCtx)),
	}

	return server
}
