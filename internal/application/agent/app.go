package agent

import (
	"app/internal/controller/agent"
	"app/internal/controller/async"
	"app/internal/dataprovider/agentapi"
	"app/internal/dataprovider/loader"
	"app/internal/dataprovider/logger"
	agentUC "app/internal/usecase/agent"
	"app/pkg/ctxtool"
	"context"
)

func Serve(ctx context.Context) {
	ctx = ctxtool.NewSystemContext(ctx, "main")

	cfg := parseFlag()

	logger := logger.New(cfg.Debug, cfg.Trace)
	logger.InfoContext(ctx, "Инициализация агента")

	async := async.New(logger)
	loader := loader.New(logger)

	agentApi := agentapi.New(agentapi.Config{
		Prefixes:     loader.Prefixes(),
		Token:        cfg.Token,
		AgentName:    cfg.Name,
		Scheme:       cfg.Scheme,
		HostWithPort: cfg.Addr,
		Logger:       logger,
	})

	useCases := agentUC.New(logger, agentApi, loader)
	controller := agent.New(logger, useCases)

	async.RegisterRunner(ctx, controller)

	logger.InfoContext(ctx, "Система запущена")

	err := async.Serve(ctx)
	if err != nil {
		logger.ErrorContext(ctx, err.Error())

		return
	}

	logger.InfoContext(ctx, "Процессы завершены, выход")
}
