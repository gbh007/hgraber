package webgin

import (
	"app/db"
	"app/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MainInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"count":               db.SelectTitlesCount(),
		"not_load_count":      db.SelectUnloadTitlesCount(),
		"page_count":          db.SelectPagesCount(),
		"not_load_page_count": db.SelectUnloadPagesCount(),
	})
}

func NewTitle(ctx *gin.Context) {
	request := struct {
		URL string `json:"url" binding:"required"`
	}{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	err = handler.FirstHandle(request.URL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	} else {
		ctx.JSON(http.StatusOK, struct{}{})
	}
}
