package webgin

import (
	"app/db"
	"app/file"
	"app/handler"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MainInfo(ctx *gin.Context) {
	sctx := getCoreContext(ctx)
	ctx.JSON(http.StatusOK, gin.H{
		"count":               db.SelectTitlesCount(sctx),
		"not_load_count":      db.SelectUnloadTitlesCount(sctx),
		"page_count":          db.SelectPagesCount(sctx),
		"not_load_page_count": db.SelectUnloadPagesCount(sctx),
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
	sctx := getCoreContext(ctx)
	err = handler.FirstHandle(sctx, request.URL)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
	} else {
		ctx.JSON(http.StatusOK, struct{}{})
	}
}

func TitleList(ctx *gin.Context) {
	request := struct {
		Count  int `json:"count" binding:"required"`
		Offset int `json:"offset"`
	}{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	sctx := getCoreContext(ctx)
	data := db.SelectTitles(sctx, request.Offset, request.Count)
	ctx.JSON(http.StatusOK, data)
}

func TitleInfo(ctx *gin.Context) {
	request := struct {
		ID int `json:"id" binding:"required"`
	}{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	sctx := getCoreContext(ctx)
	data, err := db.SelectTitleByID(sctx, request.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func TitlePage(ctx *gin.Context) {
	request := struct {
		ID   int `json:"id" binding:"required"`
		Page int `json:"page" binding:"required"`
	}{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	sctx := getCoreContext(ctx)
	data, err := db.SelectPagesByTitleIDAndNumber(sctx, request.ID, request.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, data)
}

func SaveToZIP(ctx *gin.Context) {
	request := struct {
		From int `json:"from" binding:"required"`
		To   int `json:"to" binding:"required"`
	}{}
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	sctx := getCoreContext(ctx)
	for i := request.From; i <= request.To; i++ {
		err = file.LoadToZip(sctx, i)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}
	ctx.JSON(http.StatusOK, struct{}{})
}
