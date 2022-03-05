package system

import (
	"context"
	"time"
)

var requestIDKey = 1

func NewSystemContext(parent context.Context, name string) context.Context {
	if name == "" {
		name = "SYSTEM-" + hash(time.Now().String())
	}

	return withRequestIDContext(parent, name)
}

func NewUserContext(parent context.Context) context.Context {
	return withRequestIDContext(parent, "USER-"+hash(time.Now().String()))
}

func GetRequestID(ctx context.Context) string {
	id, ok := ctx.Value(&requestIDKey).(string)
	if !ok {
		id = "???"
	}

	return id
}

func withRequestIDContext(parent context.Context, id string) context.Context {
	return &requestIDContext{
		Context:   parent,
		requestID: id,
	}
}

type requestIDContext struct {
	context.Context
	requestID string
}

func (rc *requestIDContext) Value(key interface{}) interface{} {
	if key == &requestIDKey {
		return rc.requestID
	}

	return rc.Context.Value(key)
}
