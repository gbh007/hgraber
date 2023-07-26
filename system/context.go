package system

import (
	"context"
	"errors"
	"time"
)

type contextKey struct {
	name string
}

var (
	requestIDKey               = &contextKey{"requestIDKey"}
	debugKey                   = &contextKey{"debugKey"}
	ContextAlreadyStoppedError = errors.New("ContextAlreadyStoppedError")
)

func NewSystemContext(parent context.Context, name string) context.Context {
	if name == "" {
		name = "System-" + hash(time.Now().String())
	}

	return context.WithValue(parent, requestIDKey, name)
}

func NewUserContext(parent context.Context) context.Context {
	name := "User-" + hash(time.Now().String())

	return context.WithValue(parent, requestIDKey, name)
}

func GetRequestID(ctx context.Context) string {
	id, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		id = "???"
	}

	return id
}

func IsAliveContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ContextAlreadyStoppedError
	default:
	}

	return nil
}

// WithDetach - отделяет любые операция завершения оборачивая контекст
func WithDetach(ctx context.Context) context.Context {
	return detachedContext{ctx}
}

type detachedContext struct {
	parent context.Context
}

func (_ detachedContext) Deadline() (time.Time, bool)         { return time.Time{}, false }
func (_ detachedContext) Done() <-chan struct{}               { return nil }
func (_ detachedContext) Err() error                          { return nil }
func (ctx detachedContext) Value(key interface{}) interface{} { return ctx.parent.Value(key) }
