package weather_test

import (
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"weather"
)

func TestNewClient(t *testing.T) {
	_, err := weather.NewClient("")
	if err == nil {
		t.Fatal("want error with empty API key, got nil")
	}
	APIKey := "DUMMY"
	client, err := weather.NewClient(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	if client.APIKey != APIKey {
		t.Errorf("want client.APIKey=%q, got %q", APIKey, client.APIKey)
	}
}

func TestGetWeather(t *testing.T) {
	client, err := weather.NewClient("Dummy API key")
	if err != nil {
		t.Fatal(err)
	}
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wantURL := "/data/2.5/weather"
		if r.URL.EscapedPath() != wantURL {
			t.Errorf("want URL %q, got %q", wantURL, r.URL.Path)
		}
		wantQuery := "q=London&appid=Dummy+API+key"
		if r.URL.RawQuery != wantQuery {
			t.Errorf("want query %q, got %q", wantQuery, r.URL.RawQuery)
		}
		jsonStream, err := os.Open("testdata/api_response.json")
		if err != nil {
			t.Fatal(err)
		}
		defer jsonStream.Close()
		io.Copy(w, jsonStream)
	}))
	client.APIURL = ts.URL
	client.HTTPClient = ts.Client()
	conditions, err := client.GetWeather("London")
	if err != nil {
		t.Fatal(err)
	}
	wantSummary := "Clouds"
	if wantSummary != conditions.Summary {
		t.Errorf("want summary %q, got %q", wantSummary, conditions.Summary)
	}
	wantTemp := 274.78
	if !closeEnough(wantTemp, conditions.TemperatureKelvin) {
		t.Errorf("want temperature %f, got %f", wantTemp, conditions.TemperatureKelvin)
	}
	conditions, err = client.GetWeather("London")
	if err != nil {
		t.Fatal(err)
	}
}

func closeEnough(a, b float64) bool {
	return math.Abs(a-b) < 0.00000000000001
}

func TestParseArgs(t *testing.T) {
	t.Parallel()
	tcs := []struct{
		input []string
		wantFahrMode bool
		wantLocation string
	}{
		{
			input: []string{"-fahr", "lagos"},
			wantFahrMode: true,
			wantLocation: "lagos",
		},
		{
			input: []string{"lagos"},
			wantFahrMode: false,
			wantLocation: "lagos",
		},
	}
	for _, tc := range tcs {
		fahrMode, location := weather.ParseArgs(tc.input)
		if tc.wantFahrMode != fahrMode {
			t.Errorf("input %q, want fahrMode %t, got %t", tc.input, tc.wantFahrMode, fahrMode)
		}
		if tc.wantLocation != location {
			t.Errorf("input %q, want location %q, got %q", tc.input, tc.wantLocation, location)
		}
	}
}

func TestConditionsString(t *testing.T) {
	input := weather.Conditions{
		Summary: "Sunny",
		TemperatureKelvin: 288.55,
	}
	want := "Sunny 15.4ºC"
	got := input.StringCelsius()
	if want != got {
		t.Errorf("input: %#v want %q, got %q", input, want, got)
	}
}