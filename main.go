package main

import (
	"fmt"
	"github.com/icodealot/noaa"
)

type DailyForecast struct {
	Name     string
	HighTemp float64
	LowTemp  float64
}

func sevenDayForecast(lat string, lng string) ([]DailyForecast, error) {
	forecast, err := noaa.Forecast(lat, lng)
	if err != nil {
		return nil, err
	}

	var sevenDayForecast = []DailyForecast{}
	var currentDay DailyForecast

	for i, period := range forecast.Periods {
		if i%2 == 0 {
			currentDay = DailyForecast{
				Name:     period.Name,
				HighTemp: period.Temperature,
			}
		} else {
			currentDay.LowTemp = period.Temperature
			sevenDayForecast = append(sevenDayForecast, currentDay)
		}
	}

	return sevenDayForecast, nil
}

func main() {
	forecast, err := sevenDayForecast("42.6526", "-73.7562")

	if err != nil {
		panic("Failed to retrieve forecast.")
	}

	for _, day := range forecast {
		fmt.Printf("%-20s %.0f / %.0f\n", day.Name, day.HighTemp, day.LowTemp)
	}
}
