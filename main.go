package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
	
	"github.com/hectormalot/omgo"
)

type GeoResponce struct {
	Results	[]LatLong	`json:"results"`
}

type LatLong struct {
	Latitude	float64	`json:"latitude"`
	Longitude	float64	`json:"longitude"`
}

type WeatherData struct {
	DataTime 		[]time.Time
	Temperatures	[]float64
	Precipitations	[]float64
}

func getLangLong(city string) (*LatLong, error) {
	client := http.Client{Timeout: 5 * time.Second}
	apiUlr := "https://geocoding-api.open-meteo.com/v1/search"
	req, err := http.NewRequest(http.MethodGet, apiUlr, nil)
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

func getWeather(latLong LatLong) (*WeatherData, error) {
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
	return &WeatherData{resp.HourlyTimes,
		resp.HourlyMetrics["temperature_2m"],
		resp.HourlyMetrics["precipitation_probability"]}, nil
}

func main() {
	// ...
}
