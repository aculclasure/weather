package weather_test

import (
	"os"
	"testing"

	"github.com/aculclasure/weather"
)

func TestCurrentWeatherCLI(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		apiKey      string
		args        []string
		errExpected bool
	}{
		"missing OPENWEATHER_API_KEY environment variable returns an error": {
			apiKey:      "",
			errExpected: true,
		},
		"missing weather location positional argument returns an error": {
			apiKey:      "KEY",
			args:        []string{"weathercli", "--units=imperial"},
			errExpected: true,
		},
		"missing value for units flag returns an error": {
			apiKey:      "KEY",
			args:        []string{"weathercli", "--units=", "London"},
			errExpected: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			os.Setenv("OPENWEATHER_API_KEY", tc.apiKey)
			err := weather.CurrentWeatherCLI(tc.args)
			errReceived := err != nil

			if tc.errExpected != errReceived {
				t.Fatalf("CLI(%+v) returned unexpected error status: %v", tc.args, errReceived)
			}
		})
	}
}
