package hgraberweb

import (
	"app/internal/controller/hgraberweb/internal/static"
	"app/internal/domain/hgraber"
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
)

type useCases interface {
	Info(ctx context.Context) (*hgraber.MainInfo, error)

	GetBook(ctx context.Context, id int) (hgraber.Book, error)
	GetPage(ctx context.Context, id int, page int) (*hgraber.Page, error)
	GetBooks(ctx context.Context, filter hgraber.BookFilterOuter) hgraber.FilteredBooks

	UpdatePageRate(ctx context.Context, id int, page int, rating int) error
	UpdateBookRate(ctx context.Context, id int, rating int) error

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

	logger  *slog.Logger
	webtool webtool

	addr      string
	outerAddr string
	staticDir string
	token     string
}

type Config struct {
	UseCases useCases
	Monitor  monitor

	Logger  *slog.Logger
	Webtool webtool

	Addr          string
	OuterAddr     string
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
		outerAddr: cfg.OuterAddr,
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

	mux.Handle("/api/info", tokenHandler(ws.token, ws.mainInfo()))
	mux.Handle("/api/login", ws.login(ws.token))

	mux.Handle("/api/book", tokenHandler(ws.token, ws.bookInfo()))
	mux.Handle("/api/book/new", tokenHandler(ws.token, ws.bookNew()))

	mux.Handle("/api/books", tokenHandler(ws.token, ws.bookList()))
	mux.Handle("/api/books/export", tokenHandler(ws.token, ws.booksExport()))

	mux.Handle("/api/rate", tokenHandler(ws.token, ws.ratingUpdate()))

	server := &http.Server{
		Addr: ws.addr,
		Handler: ws.webtool.PanicDefender(
			ws.webtool.CORS(mux),
		),
		BaseContext: ws.webtool.NewBaseContext(context.WithoutCancel(parentCtx)),
	}

	return server
}
