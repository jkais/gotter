package main
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Config struct {
	ApiKey string `json:"api_key"`
	City   string `json:"city"`
}

func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "gotter", "config.json")
}

func loadConfig() (*Config, error) {
	file, err := os.Open(getConfigPath())
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	return &config, err
}

type Cache struct {
	Timestamp int64  `json:"timestamp"`
	Response  string `json:"response"`
}

func getCachePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "gotter", "cache.json")
}

func ensureCacheFile() {
	cachePath := getCachePath()
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		cache := Cache{Timestamp: 0, Response: ""}
		data, _ := json.MarshalIndent(cache, "", "  ")
		os.WriteFile(cachePath, data, 0644)
	}
}

func needsUpdate() bool {
	cachePath := getCachePath()

	file, err := os.Open(cachePath)
	if err != nil {
		return true
	}
	defer file.Close()

	var cache Cache
	err = json.NewDecoder(file).Decode(&cache)
	if err != nil {
		return true
	}

	return time.Now().Unix() - cache.Timestamp > 600
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

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("Config error:", err)
		return
	}

	ensureCacheFile()

	apiKey := config.ApiKey
	city := config.City

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric", city, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("HTTP error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("API error: Status %d\n", resp.StatusCode)
		return
	}

	var data map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)

	temp := data["main"].(map[string]interface{})["temp"].(float64)
	icon := data["weather"].([]interface{})[0].(map[string]interface{})["icon"].(string)

	fmt.Printf("%.0f°C %s", temp, getWeatherIcon(icon))
}
