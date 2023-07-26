package super

import (
	"app/system"
	"context"
	"fmt"
)

type Runner interface {
	Start(context.Context) (chan struct{}, error)
	Name() string
}

type Object interface {
	Storage() Storage
	Title() TitleHandler
	RegisterRunner(context.Context, Runner)
	Run(context.Context) error
}

func NewObject(storage Storage, title TitleHandler) Object {
	return &objectImplementation{
		storage: storage,
		title:   title,
	}
}

type objectImplementation struct {
	storage        Storage
	title          TitleHandler
	runnerChannels []chan struct{}
	runners        []Runner
}

func (o *objectImplementation) Storage() Storage {
	return o.storage
}

func (o *objectImplementation) Title() TitleHandler {
	return o.title
}

func (o *objectImplementation) RegisterRunner(ctx context.Context, runner Runner) {
	o.runners = append(o.runners, runner)
}

func (o *objectImplementation) Run(parentCtx context.Context) error {
	ctx, cnl := context.WithCancel(parentCtx)
	defer cnl()

	for _, r := range o.runners {
		c, err := r.Start(ctx)
		if err != nil {
			err = fmt.Errorf("start %s: %w", r.Name(), err)

			system.Error(ctx, err)

			return err
		}

		o.runnerChannels = append(o.runnerChannels, c)
	}

	// Дожидаемся завершения потоков
	for _, c := range o.runnerChannels {
		<-c
	}

	return nil
}
