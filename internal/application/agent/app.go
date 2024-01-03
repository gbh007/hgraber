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

	logger := logger.New(false, false)
	logger.Info(ctx, "Инициализация агента")
	cfg := parseFlag()

	debug := false // FIXME: брать из конфигурации

	if debug {
		logger.SetDebug(debug)
	}

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

	logger.Info(ctx, "Система запущена")

	err := async.Serve(ctx)
	if err != nil {
		logger.Error(ctx, err)

		return
	}

	logger.Info(ctx, "Процессы завершены, выход")
}
