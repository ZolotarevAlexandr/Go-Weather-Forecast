package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


func main() {
	ginServer := gin.Default()
	ginServer.LoadHTMLGlob("views/*")
	ginServer.Static("/css", "./css/")

	ginServer.GET("/weather", serverGetHandler)
	ginServer.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	ginServer.GET("/api/weather", apiGetHandler)

	err := ginServer.Run("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
}
