package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hectormalot/omgo"
)


type LatLong struct {
	Latitude	float64	`json:"latitude"`
	Longitude	float64	`json:"longitude"`
}

type GeoResponce struct {
	Results	[]LatLong	`json:"results"`
}

type HourWeatherData struct {
	DateTime		time.Time
	Temperature		float64
	Precipitations	float64
}

type WeatherData struct {
	City		string
	Forecast	[]HourWeatherData
}


func getLangLong(city string) (*LatLong, error) {
	client := http.Client{Timeout: 5 * time.Second}
	const apiUrl = "https://geocoding-api.open-meteo.com/v1/search"

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("error while creating request: %w", err)
	}
	
	params := url.Values{}
	params.Add("name", url.QueryEscape(city))
	params.Add("count", "1")
	params.Add("languge", "en")
	params.Add("format", "json")
	req.URL.RawQuery = params.Encode()
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while getting responce: %w", err)
	}
	defer resp.Body.Close()
	
	var responce GeoResponce
	err = json.NewDecoder(resp.Body).Decode(&responce)
	if err != nil {
		return nil, fmt.Errorf("error while decoding responce: %w", err)
	}
	if len(responce.Results) == 0 {
		return nil, errors.New("no results found")
	}
	
	return &responce.Results[0], nil
}

func getWeather(latLong LatLong) ([]HourWeatherData, error) {
	client, err := omgo.NewClient()
	if err != nil {
		return nil, fmt.Errorf("error while creating omgo client: %w", err)
	}

	loc, err := omgo.NewLocation(latLong.Latitude, latLong.Longitude)
	if err != nil {
		return nil, fmt.Errorf("error while searching omgo location: %w", err)
	}

	opts := omgo.Options{
		HourlyMetrics: []string{"temperature_2m", "precipitation_probability"},
	}
	resp, err := client.Forecast(context.Background(), loc, &opts)
	if err != nil {
		return nil, fmt.Errorf("error while getting omgo responce: %w", err)
	}

	var forecast []HourWeatherData
	for i := 0; i < len(resp.HourlyTimes); i++ {
		forecast = append(forecast, HourWeatherData{resp.HourlyTimes[i],
			resp.HourlyMetrics["temperature_2m"][i],
			resp.HourlyMetrics["precipitation_probability"][i]})
	}
	return forecast, nil
}

func getWeatherData(city string) (*WeatherData, error) {
	latLong, err := getLangLong(city)
	if err != nil {
		return nil, fmt.Errorf("error while getting city coordinates: %w", err)
	}

	weather, err := getWeather(*latLong)
	if err != nil {
		return nil, fmt.Errorf("error while getting weather data: %w", err)
	}

	weatherData := WeatherData{city, weather}
	return &weatherData, nil
}

func serverGetHandler(ctx *gin.Context) {
	city := ctx.Query("city")
	weatherData, err := getWeatherData(city)
	if err != nil {
		ctx.HTML(http.StatusInternalServerError, "error.html", gin.H{"error": err.Error()})
		return
	}
	ctx.HTML(http.StatusOK, "weather.html", *weatherData)
}

func APIgetHandler(ctx *gin.Context) {
	city := ctx.Query("city")
	weatherData, err := getWeatherData(city)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, *weatherData)
}

func main() {
	ginServer := gin.Default()
	ginServer.LoadHTMLGlob("views/*")
	ginServer.Static("/css", "./css/")

	ginServer.GET("/weather", serverGetHandler)
	ginServer.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	ginServer.GET("/api/weather", APIgetHandler)

	err := ginServer.Run("localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
}
