package webgin

import (
	"app/system"

	"github.com/gin-gonic/gin"
)

const coreContextKey = "core-context-key"

func getCoreContext(ctx *gin.Context) system.Context {
	v, _ := ctx.Get(coreContextKey)
	// паника тут специально игнорируется
	return v.(system.Context)
}

func SetCoreContext(ctx *gin.Context) {
	ctx.Set(coreContextKey, system.NewUserContext())
}
