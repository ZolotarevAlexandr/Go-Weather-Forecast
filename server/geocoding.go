package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)


type LatLong struct {
	Latitude	float64	`json:"latitude"`
	Longitude	float64	`json:"longitude"`
}

type GeoResponce struct {
	Results	[]LatLong	`json:"results"`
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
