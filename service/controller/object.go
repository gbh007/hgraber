package controller

import (
	"app/system"
	"context"
	"fmt"
)

type Runner interface {
	Start(context.Context) (chan struct{}, error)
	Name() string
}

func NewObject() *Object {
	return new(Object)
}

type Object struct {
	runnerChannels []chan struct{}
	runners        []Runner
}

func (o *Object) RegisterRunner(ctx context.Context, runner Runner) {
	o.runners = append(o.runners, runner)
}

func (o *Object) Run(parentCtx context.Context) error {
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
