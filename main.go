package main
import (
	"fmt"
)

func main() {
	var weather = getCurrentWeather()

	temp := weather["main"].(map[string]interface{})["temp"].(float64)
	icon := weather["weather"].([]interface{})[0].(map[string]interface{})["icon"].(string)

	fmt.Printf("%.0f°C %s", temp, getWeatherIcon(icon))
}
