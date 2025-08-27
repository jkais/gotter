package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func getLocation() (*Location, error) {
	resp, err := http.Get("http://ip-api.com/json/?fields=lat,lon")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var loc Location
	if err := json.NewDecoder(resp.Body).Decode(&loc); err != nil {
		return nil, err
	}

	return &loc, nil
}

func getCurrentWeather() map[string]any {
	if needsUpdate() {
		var loc, err = getLocation()
		if (err != nil) {
			return getCache()
		}
		config := loadConfig()
		apiKey := config.ApiKey
		url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", loc.Lat, loc.Lon, apiKey)
		resp, err := http.Get(url)
		if (err != nil) {
			return getCache()
		}
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		setCache(string(body))
	}
	return getCache()
}
