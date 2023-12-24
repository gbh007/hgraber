package agent

import (
	"app/internal/controller/agent"
	"app/internal/controller/async"
	"app/internal/dataprovider/agentapi"
	"app/internal/dataprovider/loader"
	agentUC "app/internal/usecase/agent"
	"app/pkg/logger"
	"context"
	"fmt"
)

type App struct {
	async *async.Controller
}

func New() *App {
	return new(App)
}

func (app *App) Init(ctx context.Context) {
	cfg := parseFlag()

	logger := logger.New(false) // FIXME: брать из конфигурации
	app.async = async.New(logger)
	loader := loader.New(logger)

	agentApi := agentapi.New(agentapi.Config{
		Prefixes:     nil, // Обрабатываем все
		Token:        cfg.Token,
		AgentName:    cfg.Name,
		Scheme:       cfg.Scheme,
		HostWithPort: cfg.Addr,
		Logger:       logger,
	})

	useCases := agentUC.New(logger, agentApi, loader)
	controller := agent.New(logger, useCases)

	app.async.RegisterRunner(ctx, controller)
}

func (app *App) Serve(ctx context.Context) error {
	err := app.async.Serve(ctx)
	if err != nil {
		return fmt.Errorf("app: %w", err)
	}

	return nil
}
