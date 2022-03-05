package webgin

import (
	"app/system"
	"context"

	"github.com/gin-gonic/gin"
)

const coreContextKey = "core-context-key"

func getCoreContext(ctx *gin.Context) context.Context {
	v, _ := ctx.Get(coreContextKey)
	// паника тут специально игнорируется
	return v.(context.Context)
}

func SetCoreContext(ctx *gin.Context) {
	ctx.Set(coreContextKey, system.NewUserContext(ctx.Request.Context()))
}
