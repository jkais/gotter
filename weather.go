package main

import (
    "fmt"
    "io"
    "net/http"
)

func getCurrentWeather() map[string]interface{} {
    if needsUpdate() {
        config := loadConfig()
        apiKey := config.ApiKey
        city := config.City
        url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)
        resp, _ := http.Get(url)
        defer resp.Body.Close()
        body, _ := io.ReadAll(resp.Body)
        setCache(string(body))
    }
    return getCache()
}
