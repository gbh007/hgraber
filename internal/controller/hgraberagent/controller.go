package hgraberagent

import (
	"app/internal/domain/agent"
	"app/pkg/logger"
	"context"
	"io"
	"net"
	"net/http"
)

type useCases interface {
	UnprocessedBooks(ctx context.Context, prefixes []string, limit int) ([]agent.BookToHandle, error)
	UnprocessedPages(ctx context.Context, prefixes []string, limit int) ([]agent.PageToHandle, error)
	UpdateBook(ctx context.Context, book agent.BookToUpdate) error
	UploadPage(ctx context.Context, info agent.PageInfoToUpload, body io.Reader) error
}

type webtool interface {
	CORS(next http.Handler) http.Handler
	NewBaseContext(ctx context.Context) func(l net.Listener) context.Context
	PanicDefender(next http.Handler) http.Handler
	WriteNoContent(ctx context.Context, w http.ResponseWriter)
	WritePlain(ctx context.Context, w http.ResponseWriter, statusCode int, data string)
	ParseJSON(r *http.Request, data any) error
	WriteJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data any)
}

type Controller struct {
	logger *logger.Logger

	useCases useCases
	webtool  webtool

	addr  string
	token string
}

func New(useCases useCases, addr string, token string, logger *logger.Logger, web webtool) *Controller {
	return &Controller{
		logger:   logger,
		useCases: useCases,
		webtool:  web,

		addr:  addr,
		token: token,
	}
}

func (c *Controller) makeServer(parentCtx context.Context) *http.Server {
	mux := http.NewServeMux()

	mux.Handle(agent.EndpointBookUnprocessed, c.bookUnprocessed())
	mux.Handle(agent.EndpointBookUpdate, c.bookUpdate())
	mux.Handle(agent.EndpointPageUnprocessed, c.pageUnprocessed())
	mux.Handle(agent.EndpointPageUpload, c.pageUpload())

	server := &http.Server{
		Addr: c.addr,
		Handler: c.webtool.PanicDefender(
			c.webtool.CORS(
				c.tokenMiddleware(mux),
			),
		),
		BaseContext: func(l net.Listener) context.Context { return context.WithoutCancel(parentCtx) }, // FIXME: использовать название агента
	}

	return server
}
