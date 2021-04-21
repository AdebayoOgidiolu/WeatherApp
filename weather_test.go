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
	conditions, err := client.GetWeather("London", false)
	if err != nil {
		t.Fatal(err)
	}
	wantSummary := "Clouds"
	if wantSummary != conditions.Summary {
		t.Errorf("want summary %q, got %q", wantSummary, conditions.Summary)
	}
	wantTemp := 1.63
	if !closeEnough(wantTemp, conditions.Temperature) {
		t.Errorf("want temperature %f, got %f", wantTemp, conditions.Temperature)
	}
	conditions, err = client.GetWeather("London", true)
	if err != nil {
		t.Fatal(err)
	}
	wantFahrenheitTemp := 34.934
	if !closeEnough(wantFahrenheitTemp, conditions.Temperature) {
		t.Errorf("want temperature %f, got %f", wantFahrenheitTemp, conditions.Temperature)
	}
}

func closeEnough(a, b float64) bool {
	return math.Abs(a-b) < 0.00000000000001
}
