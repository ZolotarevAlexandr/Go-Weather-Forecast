package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"text/tabwriter"
	"time"
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

func getData(city string) (WeatherData, error) {
	client := http.Client{Timeout: 5 * time.Second}
	const apiUrl = "http://localhost:8000/api/weather"
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		return WeatherData{}, err
	}
	
	params := url.Values{}
	params.Add("city", url.QueryEscape(city))
	req.URL.RawQuery = params.Encode()
	
	resp, err := client.Do(req)
	if err != nil {
		return WeatherData{}, err
	} else if resp.StatusCode != http.StatusOK {
		return WeatherData{}, errors.New("status code not 200")
	}
	
	var result WeatherData
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return WeatherData{}, err
	}
	return result, nil
}

func main() {
	for {
		var input string
		fmt.Println("Input city name (q to exit):")
		_, err := fmt.Scanf("%s\n", &input)
		if err != nil {
			fmt.Println(err)
			continue
		}
		
		if input == "q" {
			return
		}

		forecast, err := getData(input)
		if err != nil {
			fmt.Println(err)
			continue
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		fmt.Printf("Weather for %v\n", forecast.City)
		_, err = fmt.Fprint(w, "Date\tTemperature\tPrecipitations\t\n")
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, hourForecast := range forecast.Forecast {
			_, err = fmt.Fprintf(w, "%v\t%vC°\t%v%%\t\n", hourForecast.DateTime.Format("02.01 Mon 15:04"), hourForecast.Temperature, hourForecast.Precipitations)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}

		err = w.Flush()
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
