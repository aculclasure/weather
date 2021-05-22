package weather_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aculclasure/weather"
	"github.com/google/go-cmp/cmp"
)

const nonJSONData = "123456"

func TestDecodeCurrent(t *testing.T) {
	t.Parallel()
	validData, err := ioutil.ReadFile("testdata/currentWeatherAPIResp.json")
	if err != nil {
		t.Fatal(err)
	}
	testCases := map[string]struct {
		input       []byte
		want        weather.CurrentAPIResp
		errExpected bool
	}{
		"non-json input returns an error": {
			input:       []byte(nonJSONData),
			want:        weather.CurrentAPIResp{},
			errExpected: true,
		},
		"complete json input returns CurrentAPIResp": {
			input: []byte(validData),
			want: weather.CurrentAPIResp{
				Summaries: []weather.Summary{
					{Desc: "few clouds"},
				},
				Metrics: weather.Metrics{
					Temp:     52.72,
					Humidity: 47,
				},
			},
			errExpected: false,
		},
	}

	comparer := cmp.Comparer(func(c1, c2 weather.CurrentAPIResp) bool {
		return cmp.Equal(c1.Summaries, c2.Summaries) &&
			c1.Metrics.Humidity == c2.Metrics.Humidity &&
			closeEnough(c1.Metrics.Temp, c2.Metrics.Temp)
	})

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := weather.DecodeCurrent(tc.input)
			errReceived := err != nil

			if tc.errExpected != errReceived {
				t.Fatalf("got unexpected error status: %v", errReceived)
			}

			if !tc.errExpected && !cmp.Equal(tc.want, got, comparer) {
				t.Fatalf("want != got\ndiff=%s", cmp.Diff(tc.want, got, comparer))
			}
		})
	}
}

func TestGetCurrentWeatherData(t *testing.T) {
	t.Parallel()
	validData, err := ioutil.ReadFile("testdata/currentWeatherAPIResp.json")
	if err != nil {
		t.Fatal(err)
	}
	client, err := weather.NewClient("apikey")
	if err != nil {
		t.Fatalf("got error creating new weather client: %v", err)
	}

	wantReqURI := "/data/2.5/weather?q=London&units=imperial&appid=apikey"
	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wantReqURI != r.RequestURI {
			t.Fatalf("want request URI: %s, got %s", wantReqURI, r.RequestURI)
		}
		fmt.Fprint(w, string(validData))
	}))
	defer testServer.Close()
	client.HTTPClient = testServer.Client()
	client.BaseURL = testServer.URL
	gotData, err := client.Current("London", "imperial")
	if err != nil {
		t.Fatal(err)
	}
	wantData := validData
	if !bytes.Equal(wantData, gotData) {
		t.Fatalf("want != got\ndiff=%s", cmp.Diff(wantData, gotData))
	}
}

func TestGetCurrentWeatherWithInvalidArgumentsReturnsError(t *testing.T) {
	t.Parallel()
	testCases := map[string]struct {
		location string
		units    string
	}{
		"Empty location": {
			location: "",
			units:    "standard",
		},
		"Invalid units": {
			location: "new york city,ny,us",
			units:    "not a unit",
		},
	}
	client, err := weather.NewClient("apikey")
	if err != nil {
		t.Fatalf("got error creating weather.Client: %v", err)
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_, err := client.Current(tc.location, tc.location)
			if err == nil {
				t.Fatalf("client.Current(%s, %s) did not return an expected error",
					tc.location, tc.units)
			}
		})
	}
}

func TestDecodeGeoData(t *testing.T) {
	t.Parallel()
	emptyData := "[]"
	validData, err := ioutil.ReadFile("testdata/geocodeAPIResp.json")
	if err != nil {
		t.Fatalf("unable to read test data file: %v", err)
	}
	testCases := map[string]struct {
		input       []byte
		want        weather.Location
		errExpected bool
	}{
		"non-json input returns an error": {
			input:       []byte(nonJSONData),
			errExpected: true,
		},
		"empty response returns an error": {
			input:       []byte(emptyData),
			errExpected: true,
		},
		"valid geocode data returns a weather.GeoCodeSnippet": {
			input: validData,
			want: weather.Location{
				Name:    "London",
				Country: "GB",
				Lat:     51.5085,
				Lon:     -0.1257,
			},
			errExpected: false,
		},
	}
	comparer := cmp.Comparer(func(loc1, loc2 weather.Location) bool {
		return loc1.Country == loc2.Country &&
			loc1.Name == loc2.Name &&
			closeEnough(loc1.Lat, loc2.Lat) &&
			closeEnough(loc1.Lon, loc2.Lon)
	})

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := weather.DecodeGeoData(tc.input)
			errReceived := err != nil

			if tc.errExpected != errReceived {
				t.Fatalf("got unexpected error status: %v", errReceived)
			}

			if !tc.errExpected && !cmp.Equal(tc.want, got, comparer) {
				t.Fatalf("want != got\ndiff=%s", cmp.Diff(tc.want, got, comparer))
			}
		})
	}
}

func TestGetGeocodeData(t *testing.T) {
	t.Parallel()
	client, err := weather.NewClient("apikey")
	if err != nil {
		t.Fatalf("got error creating new weather client: %v", err)
	}
	wantData, err := ioutil.ReadFile("testdata/geocodeAPIResp.json")
	if err != nil {
		t.Fatalf("unable to read test data file: %v", err)
	}
	wantReqURI := "/geo/1.0/direct?q=London,GB&limit=1&appid=apikey"
	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wantReqURI != r.RequestURI {
			t.Fatalf("want request URI: %s, got %s", wantReqURI, r.RequestURI)
		}
		fmt.Fprint(w, string(wantData))
	}))
	defer testServer.Close()
	client.HTTPClient = testServer.Client()
	client.BaseURL = testServer.URL
	gotData, err := client.GeocodeData("London,GB")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(wantData, gotData) {
		t.Fatalf("want != got\ndiff=%s", cmp.Diff(wantData, gotData))
	}
}

func TestGetGeocodeDataWithoutLocationReturnsError(t *testing.T) {
	client, err := weather.NewClient("apikey")
	if err != nil {
		t.Fatalf("got error creating new weather client: %v", err)
	}
	invalidLoc := ""
	_, err = client.GeocodeData(invalidLoc)
	if err == nil {
		t.Fatalf("client.Current(%s) did not return an expected error", invalidLoc)
	}
}

func TestDecodeOneCallDailyData(t *testing.T) {
	t.Parallel()
	t.Run("Empty data slice argument returns an error", func(t *testing.T) {
		emptyData := ""
		_, err := weather.DecodeOneCallDailyData([]byte(emptyData))
		if err == nil {
			t.Fatalf("wanted an error but did not get one")
		}
	})

	t.Run("Valid data gets decoded into []OneCallDayForecast", func(t *testing.T) {
		validData, err := ioutil.ReadFile("testdata/oneCallAPIResp.json")
		if err != nil {
			t.Fatalf("unable to read test data file: %v", err)
		}
		dayForecasts, err := weather.DecodeOneCallDailyData(validData)
		if err != nil {
			t.Fatalf("DecodeOneCallDailyData(data) returned unexpected error %v\nfor data:\n%s",
				err, string(validData))
		}

		wantLength := 8
		gotLength := len(dayForecasts)
		if wantLength != gotLength {
			t.Fatalf("want []OneCallDayForecast to have length %d, got %d",
				wantLength, gotLength)
		}

		comparer := cmp.Comparer(func(f1, f2 weather.OneCallDayForecast) bool {
			return f1.Date == f2.Date &&
				closeEnough(f1.Temp.Low, f2.Temp.Low) &&
				closeEnough(f1.Temp.High, f2.Temp.High) &&
				f1.Humidity == f2.Humidity &&
				cmp.Equal(f1.Weather, f2.Weather)
		})
		wantFirstDayForecast := weather.OneCallDayForecast{
			Date:     1621360800,
			Temp:     weather.OneCallDayTemp{Low: 290.44, High: 298.72},
			Humidity: 72,
			Weather:  []weather.OneCallDaySummary{{Desc: "very heavy rain"}},
		}
		gotFirstDayForecast := dayForecasts[0]
		if !cmp.Equal(wantFirstDayForecast, gotFirstDayForecast, comparer) {
			t.Fatalf("want != got\ndiff=%s", cmp.Diff(
				wantFirstDayForecast,
				gotFirstDayForecast,
				comparer,
			))
		}
	})
}

func TestGetOneCallData(t *testing.T) {
	t.Parallel()
	client, err := weather.NewClient("apikey")
	if err != nil {
		t.Fatalf("got error creating new weather client: %v", err)
	}
	wantData, err := ioutil.ReadFile("testdata/oneCallAPIResp.json")
	if err != nil {
		t.Fatalf("unable to read test data file: %v", err)
	}
	wantReqURI := "/data/2.5/onecall?lat=33.44&lon=-94.04&units=standard&appid=apikey"
	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wantReqURI != r.RequestURI {
			t.Fatalf("want request URI: %s, got %s", wantReqURI, r.RequestURI)
		}
		fmt.Fprint(w, string(wantData))
	}))
	defer testServer.Close()
	client.HTTPClient = testServer.Client()
	client.BaseURL = testServer.URL
	gotData, err := client.OneCallData(33.44, -94.04, "standard")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(wantData, gotData) {
		t.Fatalf("want != got\ndiff=%s", cmp.Diff(wantData, gotData))
	}
}

func TestGetOneCallDataEncodesExcludedTimeFramesInRequest(t *testing.T) {
	t.Parallel()
	client, err := weather.NewClient("apikey")
	if err != nil {
		t.Fatalf("got error creating new weather client: %v", err)
	}
	wantReqURI := "/data/2.5/onecall?lat=33.44&lon=-94.04&units=standard&appid=apikey&exclude=current,minutely,hourly,alerts"
	testServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if wantReqURI != r.RequestURI {
			t.Fatalf("want request URI: %s, got %s", wantReqURI, r.RequestURI)
		}
		fmt.Fprint(w, "ok")
	}))
	defer testServer.Close()
	client.HTTPClient = testServer.Client()
	client.BaseURL = testServer.URL
	_, err = client.OneCallData(
		33.44,
		-94.04,
		"standard",
		[]string{"ignored", "current", "ignored", "minutely", "hourly", "alerts", "yearly-ignored"}...)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetOneCallDataWithInvalidUnitsReturnsError(t *testing.T) {
	t.Parallel()
	client, err := weather.NewClient("apikey")
	if err != nil {
		t.Fatalf("got error creating new weather client: %v", err)
	}
	lat, lon, invalidUnits := 1.00, 2.00, "martian"
	_, err = client.OneCallData(lat, lon, invalidUnits)
	if err == nil {
		t.Fatalf("OneCallData(%.2f, %.2f, %s) did not return an expected error",
			lat, lon, invalidUnits)
	}
}

// closeEnough accepts 2 float64 values and returns true if they are
// within 0.001 of each other or returns false otherwise.
func closeEnough(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}
