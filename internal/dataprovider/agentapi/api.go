package agentapi

import (
	"log/slog"
	"net/http"
	"time"
)

const (
	apiTimeout = time.Second * 10
	apiName    = "agent"
)

type Config struct {
	Prefixes  []string
	Token     string
	AgentName string

	Scheme       string
	HostWithPort string

	Logger *slog.Logger
}

type API struct {
	prefixes  []string
	token     string
	agentName string

	scheme       string
	hostWithPort string

	client *http.Client

	logger *slog.Logger
}

func New(cfg Config) *API {
	return &API{
		prefixes:  cfg.Prefixes,
		token:     cfg.Token,
		agentName: cfg.AgentName,

		scheme:       cfg.Scheme,
		hostWithPort: cfg.HostWithPort,

		logger: cfg.Logger,

		client: &http.Client{
			Timeout: apiTimeout,
		},
	}
}
