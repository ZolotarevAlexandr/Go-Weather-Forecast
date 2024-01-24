package main

import (
	"context"
	"fmt"
	"time"

	"github.com/hectormalot/omgo"
)


type HourWeatherData struct {
	DateTime		time.Time
	Temperature		float64
	Precipitations	float64
}

type WeatherData struct {
	City		string
	Forecast	[]HourWeatherData
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
