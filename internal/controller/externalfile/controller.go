package externalfile

import (
	"app/internal/dto"
	"app/pkg/webtool"
	"app/system"
	"context"
	"io"
	"net/http"
)

type fileStorage interface {
	CreatePageFile(ctx context.Context, id, page int, ext string) (io.WriteCloser, error)
	OpenPageFile(ctx context.Context, id, page int, ext string) (io.ReadCloser, error)
	CreateExportFile(ctx context.Context, name string) (io.WriteCloser, error)
}

type Controller struct {
	fileStorage fileStorage

	addr  string
	token string
}

func New(fileStorage fileStorage, addr string, token string) *Controller {
	return &Controller{
		fileStorage: fileStorage,

		addr:  addr,
		token: token,
	}
}

func (c *Controller) makeServer(parentCtx context.Context) *http.Server {
	mux := http.NewServeMux()

	mux.Handle(dto.ExternalFileEndpointPage, c.pageHandler())
	mux.Handle(dto.ExternalFileEndpointExport, c.fileExport())

	server := &http.Server{
		Addr: c.addr,
		Handler: webtool.PanicDefender(
			webtool.Stopwatch(
				webtool.CORS(
					c.tokenMiddleware(mux),
				),
			),
		),
		ErrorLog:    system.StdErrorLogger(parentCtx),
		BaseContext: webtool.NewBaseContext(context.WithoutCancel(parentCtx)),
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
			webtool.WritePlain(r.Context(), w, http.StatusMethodNotAllowed, "unsupported method")
		}
	})
}
