package base

import (
	"app/system"
	"context"
	"net"
	"net/http"
)

type contextKey struct {
	name string
}

// ключи для ResponseContext
var (
	responseDataKey  = &contextKey{name: "responseDataKey"}
	responseErrorKey = &contextKey{name: "responseErrorKey"}
)

// ResponseContext контекст для ответа
type ResponseContext struct {
	context.Context
	key interface{}
	val interface{}
}

func (rc *ResponseContext) Value(key interface{}) interface{} {
	if key == rc.key {
		return rc.val
	}
	return rc.Context.Value(key)
}

func withResponseContext(parent context.Context, key, val interface{}) context.Context {
	return &ResponseContext{
		Context: parent,
		key:     key,
		val:     val,
	}
}

func SetError(r *http.Request, err error) {
	*r = *r.WithContext(
		withResponseContext(
			r.Context(),
			responseErrorKey,
			err,
		),
	)
}

func SetResponse(r *http.Request, data interface{}) {
	*r = *r.WithContext(
		withResponseContext(
			r.Context(),
			responseDataKey,
			data,
		),
	)
}

func GetError(ctx context.Context) error {
	tmp := ctx.Value(responseErrorKey)
	if tmp != nil {
		err, _ := tmp.(error)
		return err
	}

	return nil
}

func GetResponse(ctx context.Context) interface{} {
	return ctx.Value(responseDataKey)
}

func NewBaseContext(ctx context.Context) func(l net.Listener) context.Context {
	return func(l net.Listener) context.Context { return system.NewUserContext(ctx) }
}
