package webgin

import (
	"app/system/coreContext"

	"github.com/gin-gonic/gin"
)

const coreContextKey = "core-context-key"

func getCoreContext(ctx *gin.Context) coreContext.CoreContext {
	v, _ := ctx.Get(coreContextKey)
	// паника тут специально игнорируется
	return v.(coreContext.CoreContext)
}

func SetCoreContext(ctx *gin.Context) {
	ctx.Set(coreContextKey,coreContext.NewUserContext())
}
