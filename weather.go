package main

import (
    "fmt"
    "io"
    "net/http"
)

func getWeatherIcon(iconCode string) string {
    icons := map[string]string{
        "01d": "☀️", "01n": "🌙",
        "02d": "⛅", "02n": "☁️",
        "03d": "☁️", "03n": "☁️",
        "04d": "☁️", "04n": "☁️",
        "09d": "🌦️", "09n": "🌧️",
        "10d": "🌦️", "10n": "🌧️",
        "11d": "⛈️", "11n": "⛈️",
        "13d": "🌨️", "13n": "❄️",
        "50d": "🌫️", "50n": "🌫️",
    }
    if emoji, exists := icons[iconCode]; exists {
        return emoji
    }
    return "🌡️"
}

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
