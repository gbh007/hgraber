package hgraberweb

import (
	"app/internal/controller/hgraberweb/internal/static"
	"app/internal/domain/hgraber"
	"context"
	"io"
	"net"
	"net/http"
)

type logger interface {
	Error(ctx context.Context, err error)
	IfErr(ctx context.Context, err error)
	IfErrFunc(ctx context.Context, f func() error)
	Info(ctx context.Context, args ...any)
}

type useCases interface {
	Info(ctx context.Context) (*hgraber.MainInfo, error)

	GetBook(ctx context.Context, id int) (hgraber.Book, error)
	GetPage(ctx context.Context, id int, page int) (*hgraber.Page, error)
	GetBooks(ctx context.Context, filter hgraber.BookFilter) []hgraber.Book

	UpdatePageRate(ctx context.Context, id int, page int, rate int) error
	UpdateBookRate(ctx context.Context, id int, rate int) error

	ExportBooksToZip(ctx context.Context, from, to int) error

	FirstHandle(ctx context.Context, u string) error
	FirstHandleMultiple(ctx context.Context, data []string) (*hgraber.FirstHandleMultipleResult, error)

	PageWithBody(ctx context.Context, bookID int, pageNumber int) (*hgraber.Page, io.ReadCloser, error)
}

type webtool interface {
	CORS(next http.Handler) http.Handler
	NewBaseContext(ctx context.Context) func(l net.Listener) context.Context
	PanicDefender(next http.Handler) http.Handler
	ParseJSON(r *http.Request, data any) error
	WriteJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data any)
	WriteNoContent(ctx context.Context, w http.ResponseWriter)
	WritePlain(ctx context.Context, w http.ResponseWriter, statusCode int, data string)
}

type monitor interface {
	Info() []hgraber.MonitorStat
}

type WebServer struct {
	useCases useCases
	monitor  monitor

	logger  logger
	webtool webtool

	addr      string
	outerAddr string
	staticDir string
	token     string
}

type Config struct {
	UseCases useCases
	Monitor  monitor

	Logger  logger
	Webtool webtool

	Addr          string
	Token         string
	StaticDirPath string
}

func New(cfg Config) *WebServer {
	return &WebServer{
		useCases: cfg.UseCases,
		monitor:  cfg.Monitor,

		logger:  cfg.Logger,
		webtool: cfg.Webtool,

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
		Handler: ws.webtool.PanicDefender(
			ws.webtool.CORS(mux),
		),
		BaseContext: ws.webtool.NewBaseContext(context.WithoutCancel(parentCtx)),
	}

	return server
}
