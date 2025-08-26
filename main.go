package main
import (
	"fmt"
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

func main() {
	var weather = getCurrentWeather()

	temp := weather["main"].(map[string]interface{})["temp"].(float64)
	icon := weather["weather"].([]interface{})[0].(map[string]interface{})["icon"].(string)

	fmt.Printf("%.0f°C %s", temp, getWeatherIcon(icon))
}
