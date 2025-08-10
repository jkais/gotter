package main
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type Config struct {
	ApiKey string `json:"api_key"`
	City   string `json:"city"`
}

func loadConfig() (*Config, error) {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".config", "gotter", "config.json")

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	err = json.NewDecoder(file).Decode(&config)
	return &config, err
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
