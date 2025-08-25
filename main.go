package main
import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	ApiKey string `json:"api_key"`
	City   string `json:"city"`
}

type Cache struct {
	Timestamp int64  `json:"timestamp"`
	Response  map[string]interface{} `json:"response"`
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "gotter", "config.json")
}

func loadConfig() *Config {
	file, _ := os.Open(getConfigPath())
	defer file.Close()

	var config Config
	json.NewDecoder(file).Decode(&config)
	return &config
}

func getCachePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "gotter", "cache.json")
}

func ensureCacheFile() {
	cachePath := getCachePath()
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		cache := Cache{Timestamp: 0, Response: make(map[string]interface{})}
		data, _ := json.MarshalIndent(cache, "", "  ")
		os.WriteFile(cachePath, data, 0644)
	}
}

func getCacheData() Cache {
	ensureCacheFile()
	cachePath := getCachePath()

	file, _ := os.Open(cachePath)
	defer file.Close()

	var cache Cache
	json.NewDecoder(file).Decode(&cache)
	return cache
}

func needsUpdate() bool {
	return time.Now().Unix() - getCacheData().Timestamp > 600
}

func getWeatherIcon(iconCode string) string {
	icons := map[string]string{
		"01d": "☀️", "01n": "🌙",  // clear sky
		"02d": "⛅", "02n": "☁️",  // few clouds  
		"03d": "☁️", "03n": "☁️",  // scattered clouds
		"04d": "☁️", "04n": "☁️",  // broken clouds
		"09d": "🌦️", "09n": "🌧️", // shower rain
		"10d": "🌦️", "10n": "🌧️", // rain
		"11d": "⛈️", "11n": "⛈️",  // thunderstorm
		"13d": "🌨️", "13n": "❄️",  // snow
		"50d": "🌫️", "50n": "🌫️", // mist
	}
	if emoji, exists := icons[iconCode]; exists {
		return emoji
	}
	return "🌡️" // fallback
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

func setCache(jsonResponse string) {
	var responseData map[string]interface{}
	json.Unmarshal([]byte(jsonResponse), &responseData)
	prettyResponse, _ := json.MarshalIndent(responseData, "  ", "  ")

	cache := fmt.Sprintf(`{
		"timestamp": %d,
		"response": %s
		}`, time.Now().Unix(), string(prettyResponse))

	os.WriteFile(getCachePath(), []byte(cache), 0644)
}

func getCache() map[string]interface{} {
	cache := getCacheData()
	return cache.Response
}

func main() {
	var weather = getCurrentWeather()

	temp := weather["main"].(map[string]interface{})["temp"].(float64)
	icon := weather["weather"].([]interface{})[0].(map[string]interface{})["icon"].(string)

	fmt.Printf("%.0f°C %s", temp, getWeatherIcon(icon))
}
