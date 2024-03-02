package slogHandler

type handlerConfig struct {
	isDebug bool
	hooks   []handlerCtxHook
	printer Printer
}

type handlerOption interface {
	apply(*handlerConfig)
}

type applyer func(*handlerConfig)

func (a applyer) apply(cfg *handlerConfig) {
	a(cfg)
}

func WithDebug() applyer {
	return applyer(func(hc *handlerConfig) {
		hc.isDebug = true
	})
}

func WithCtxHooks(hooks ...handlerCtxHook) applyer {
	return applyer(func(hc *handlerConfig) {
		hc.hooks = append(hc.hooks, hooks...)
	})
}

func WithPrinter(printer Printer) applyer {
	return applyer(func(hc *handlerConfig) {
		hc.printer = printer
	})
}
