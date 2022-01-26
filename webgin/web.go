package webgin

import (
	"app/config"
	"log"

	"github.com/gin-gonic/gin"
)

// Run запускает веб сервер
func Run(addr string) <-chan struct{} {
	router := gin.New()
	// обработчик статики
	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/list", "./static/list.html")
	router.StaticFile("/read", "./static/read.html")
	router.Static("/static", "./static")
	router.Static("/file", config.DefaultFilePath)

	router.GET("/info", MainInfo)
	router.POST("/new", NewTitle)
	router.POST("/title/list", TitleList)
	router.POST("/title/details", TitleInfo)
	router.POST("/title/page", TitlePage)

	// mux.HandleFunc("/prepare", SaveToZIP)

	done := make(chan struct{})
	go func() {
		if err := router.Run(addr); err != nil {
			log.Println(err)
		}
		close(done)
	}()
	return done
}
