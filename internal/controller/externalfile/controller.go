package externalfile

import (
	"app/internal/domain/externalfile"
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
)

type fileStorage interface {
	CreatePageFile(ctx context.Context, id, page int, ext string, body io.Reader) error
	OpenPageFile(ctx context.Context, id, page int, ext string) (io.ReadCloser, error)
	CreateExportFile(ctx context.Context, name string, body io.Reader) error
}

type webtool interface {
	CORS(next http.Handler) http.Handler
	NewBaseContext(ctx context.Context) func(l net.Listener) context.Context
	PanicDefender(next http.Handler) http.Handler
	WriteNoContent(ctx context.Context, w http.ResponseWriter)
	WritePlain(ctx context.Context, w http.ResponseWriter, statusCode int, data string)
	MethodSplitter(handlers map[string]http.Handler) http.Handler
}

type Controller struct {
	logger *slog.Logger

	fileStorage fileStorage
	webtool     webtool

	addr  string
	token string
}

func New(fileStorage fileStorage, addr string, token string, logger *slog.Logger, web webtool) *Controller {
	return &Controller{
		logger:      logger,
		fileStorage: fileStorage,
		webtool:     web,

		addr:  addr,
		token: token,
	}
}

func (c *Controller) makeServer(parentCtx context.Context) *http.Server {
	mux := http.NewServeMux()

	mux.Handle(externalfile.EndpointPage, c.webtool.MethodSplitter(
		map[string]http.Handler{
			http.MethodGet:  c.getPage(),
			http.MethodPost: c.setPage(),
		},
	))
	mux.Handle(externalfile.EndpointExport, c.fileExport())

	server := &http.Server{
		Addr: c.addr,
		Handler: c.webtool.PanicDefender(
			c.webtool.CORS(
				c.tokenMiddleware(mux),
			),
		),
		BaseContext: c.webtool.NewBaseContext(context.WithoutCancel(parentCtx)),
	}

	return server
}
