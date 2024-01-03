package ctxtool

import (
	"context"
	"crypto/md5"
	"fmt"
	"time"
)

type contextKey struct {
	name string
}

var (
	requestIDKey = &contextKey{"requestIDKey"}
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

func hash(s string) string { return fmt.Sprintf("%x", md5.Sum([]byte(s)))[:6] }
