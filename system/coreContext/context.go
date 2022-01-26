package coreContext

import (
	"context"
	"time"
)

type CoreContext interface {
	context.Context
	SetRequestID(id string)
	GetRequestID() string
	Created() time.Time
	IsSystem() bool
}

func NewSystemContext() CoreContext {
	return &sContext{
		Context:   context.Background(),
		requestID: "SYSTEM-" + hash(time.Now().String()),
		isSystem:  true,
	}
}

func NewUserContext() CoreContext {
	return &sContext{
		Context:   context.Background(),
		requestID: "USER-" + hash(time.Now().String()),
	}
}

type sContext struct {
	context.Context
	requestID string
	created   time.Time
	isSystem  bool
}

func (sc *sContext) SetRequestID(id string) { sc.requestID = id }
func (sc sContext) GetRequestID() string    { return sc.requestID }
func (sc sContext) Created() time.Time      { return sc.created }
func (sc sContext) IsSystem() bool          { return sc.isSystem }
