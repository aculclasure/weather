// Package weather contains types and functions for communicating with
// the OpenWeatherMap API.
package weather

import (
	"fmt"
	"strings"
)

var temperatureInitials = map[string]string{
	"standard": "K",
	"metric":   "C",
	"imperial": "F",
}

// Conditions accepts a location (e.g. "london", "tampa,us", etc.), a
// measurement unit for describing weather metrics (e.g. "metric",
// "standard", "imperial"), and an OpenWeatherMap API key, makes a
// request to the OpenWeatherMap current weather API and returns a
// string summarizing the current weather for that location. An
// error is returned if the Client struct cannot be created, if
// the request to the OpenWeatherMap current weather API fails, or
// if the API response cannot be decoded properly.
func Conditions(location, units, apiKey string) (string, error) {
	client, err := NewClient(apiKey)
	if err != nil {
		return "", err
	}
	data, err := client.Current(location, units)
	if err != nil {
		return "", err
	}
	resp, err := DecodeCurrent(data)
	if err != nil {
		return "", err
	}
	desc := ""
	if len(resp.Summaries) > 0 {
		desc = resp.Summaries[0].Desc + " "
	}
	ti := temperatureInitials[units]
	return fmt.Sprintf("%s, %.2f %s, humidity %d%%",
		strings.TrimSpace(desc),
		resp.Metrics.Temp, ti,
		resp.Metrics.Humidity), nil
}
