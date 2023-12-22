package webServer

import (
	"app/internal/domain"
	"app/internal/service/webServer/static"
	"app/pkg/logger"
	"app/pkg/webtool"
	"context"
	"io"
	"net/http"
)

type pageHandler interface {
	ExportBooksToZip(ctx context.Context, from, to int) error
}

type titleHandler interface {
	// FirstHandle обрабатывает данные тайтла (новое добавление, упрощенное без парса страниц)
	FirstHandle(ctx context.Context, u string) error
	FirstHandleMultiple(ctx context.Context, data []string) domain.FirstHandleMultipleResult
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

	logger *logger.Logger

	addr      string
	outerAddr string
	staticDir string
	token     string
}

type Config struct {
	Storage storage
	Book    titleHandler
	Page    pageHandler
	Files   files
	Monitor monitor

	Logger *logger.Logger

	Addr          string
	Token         string
	StaticDirPath string
}

func New(cfg Config) *WebServer {
	return &WebServer{
		storage: cfg.Storage,
		title:   cfg.Book,
		page:    cfg.Page,
		files:   cfg.Files,
		monitor: cfg.Monitor,

		logger: cfg.Logger,

		addr:      cfg.Addr,
		outerAddr: "http://" + cfg.Addr,
		staticDir: cfg.StaticDirPath,
		token:     cfg.Token,
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
	mux.Handle("/file/", tokenHandler(ws.token, http.StripPrefix("/file/", ws.getFile())))

	// API
	mux.Handle("/auth/login", ws.routeLogin(ws.token))
	mux.Handle("/info", tokenHandler(ws.token, ws.routeMainInfo()))
	mux.Handle("/new", tokenHandler(ws.token, ws.routeNewTitle()))
	mux.Handle("/title/list", tokenHandler(ws.token, ws.routeTitleList()))
	mux.Handle("/title/details", tokenHandler(ws.token, ws.routeTitleInfo()))
	mux.Handle("/title/page", tokenHandler(ws.token, ws.routeTitlePage()))
	mux.Handle("/to-zip", tokenHandler(ws.token, ws.routeSaveToZIP()))
	mux.Handle("/app/info", tokenHandler(ws.token, ws.routeAppInfo()))
	mux.Handle("/title/rate", tokenHandler(ws.token, ws.routeSetTitleRate()))
	mux.Handle("/title/page/rate", tokenHandler(ws.token, ws.routeSetPageRate()))

	server := &http.Server{
		Addr: ws.addr,
		Handler: webtool.PanicDefender(
			webtool.CORS(mux),
		),
		BaseContext: webtool.NewBaseContext(context.WithoutCancel(parentCtx)),
	}

	return server
}
