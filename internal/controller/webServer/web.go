package webServer

import (
	"app/internal/controller/webServer/internal/static"
	"app/internal/domain"
	"app/pkg/logger"
	"app/pkg/webtool"
	"context"
	"io"
	"net/http"
)

type useCases interface {
	Info(ctx context.Context) (*domain.MainInfo, error)

	GetBook(ctx context.Context, id int) (domain.Book, error)
	GetPage(ctx context.Context, id int, page int) (*domain.Page, error)
	GetBooks(ctx context.Context, filter domain.BookFilter) []domain.Book

	UpdatePageRate(ctx context.Context, id int, page int, rate int) error
	UpdateBookRate(ctx context.Context, id int, rate int) error

	ExportBooksToZip(ctx context.Context, from, to int) error

	FirstHandle(ctx context.Context, u string) error
	FirstHandleMultiple(ctx context.Context, data []string) (*domain.FirstHandleMultipleResult, error)

	PageWithBody(ctx context.Context, bookID int, pageNumber int) (*domain.Page, io.ReadCloser, error)
}

type monitor interface {
	Info() []domain.MonitorStat
}

type WebServer struct {
	useCases useCases
	monitor  monitor

	logger *logger.Logger

	addr      string
	outerAddr string
	staticDir string
	token     string
}

type Config struct {
	UseCases useCases
	Monitor  monitor

	Logger *logger.Logger

	Addr          string
	Token         string
	StaticDirPath string
}

func New(cfg Config) *WebServer {
	return &WebServer{
		useCases: cfg.UseCases,
		monitor:  cfg.Monitor,

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
