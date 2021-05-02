package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Client represents an OpenWeatherMap API client.
type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	APIKey     string
}

// NewClient accepts an OpenWeatherMap API key as a string, creates a Client
// for communicating with the OpenWeatherMap API(s) and returns it. An error
// is returned if the apiKey argument is empty.
func NewClient(apiKey string) (Client, error) {
	if apiKey == "" {
		return Client{}, errors.New("apiKey argument must not be empty")
	}

	hc := http.DefaultClient
	hc.Timeout = 10 * time.Second
	return Client{
		HTTPClient: hc,
		BaseURL:    "https://api.openweathermap.org",
		APIKey:     apiKey,
	}, nil
}

// Current accepts a location (e.g. "london", "tampa,us", etc.), a measurement
// unit ("standard", "metric", or "imperial"), makes a call to the
// OpenWeatherMap Current Weather API to retrieve the current weather
// data for that location and returns the API response as a slice of bytes.
// An error is returned if the location or units arguments are invalid, if
// the HTTP request to the OpenWeatherMap API fails, or if there is a problem
// reading the response body.
func (c Client) Current(location, units string) ([]byte, error) {
	if location == "" {
		return nil, errors.New("location argument must not be empty")
	}
	if units != "standard" && units != "metric" && units != "imperial" {
		return nil, errors.New("units must be one of: standard, metric, imperial")
	}

	URL := fmt.Sprintf("%s/data/2.5/weather?q=%s&units=%s&appid=%s", c.BaseURL, location, units, c.APIKey)
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return nil, fmt.Errorf("error getting data from %s: %v", URL, err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("got error reading response body: %v", err)
	}

	return data, nil
}

// CurrentAPIResp represents a response from a call to the current weather
// API at OpenWeather.
type CurrentAPIResp struct {
	Summaries []Summary `json:"weather"`
	Metrics   Metrics   `json:"main"`
}

// Summary represents a weather description, like "drizzly", "overcast", etc.
type Summary struct {
	Desc string `json:"description"`
}

// Metrics represents a type to store weather metrics.
type Metrics struct {
	Temp     float64 `json:"temp"`
	Humidity int     `json:"humidity"`
}

// DecodeCurrent accepts a slice of bytes containing the response from a call
// to the OpenWeather Current Weather API, attempts to decode it into a
// Snippet, and returns the Snippet. An error is returned if
// the decoding fails.
func DecodeCurrent(data []byte) (CurrentAPIResp, error) {
	var resp CurrentAPIResp

	if err := json.Unmarshal(data, &resp); err != nil {
		return CurrentAPIResp{},
			fmt.Errorf("got error unmarshaling json %+v: %v", data, err)
	}

	return resp, nil
}
