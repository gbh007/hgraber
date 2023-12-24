package externalfile

import (
	"app/internal/domain/externalfile"
	"app/pkg/logger"
	"context"
	"io"
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
}

type Controller struct {
	logger *logger.Logger

	fileStorage fileStorage
	webtool     webtool

	addr  string
	token string
}

func New(fileStorage fileStorage, addr string, token string, logger *logger.Logger, web webtool) *Controller {
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

	mux.Handle(externalfile.EndpointPage, c.pageHandler())
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

func (c *Controller) pageHandler() http.Handler {
	getPage := c.getPage()
	setPage := c.setPage()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getPage.ServeHTTP(w, r)

		case http.MethodPost:
			setPage.ServeHTTP(w, r)

		default:
			c.webtool.WritePlain(r.Context(), w, http.StatusMethodNotAllowed, "unsupported method")
		}
	})
}
