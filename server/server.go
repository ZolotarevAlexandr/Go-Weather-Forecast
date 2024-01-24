package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


func serverGetHandler(ctx *gin.Context) {
	city := ctx.Query("city")
	weatherData, err := getWeatherData(city)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	ctx.HTML(http.StatusOK, "weather.html", *weatherData)
}

func apiGetHandler(ctx *gin.Context) {
	city := ctx.Query("city")
	weatherData, err := getWeatherData(city)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, *weatherData)
}

