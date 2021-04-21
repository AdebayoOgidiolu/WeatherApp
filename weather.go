package weather

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type APIResponse struct {
	Weather []struct {
		Main string
	}
	Main struct {
		Temp float64
	}
}

type Conditions struct {
	Summary     string
	TemperatureKelvin float64
}

func (c Conditions) StringCelsius() string {
	temp := c.TemperatureKelvin-273.15
	return fmt.Sprintf("%s %.1fºC", c.Summary, temp)
}

func (c Conditions) StringFahrenheit() string {
	temp := (c.TemperatureKelvin-273.15)*9/5 + 32
	return fmt.Sprintf("%s %.1fºF", c.Summary, temp)
}

type Client struct {
	APIKey     string
	APIURL     string
	HTTPClient *http.Client
}

func NewClient(APIKey string) (*Client, error) {
	if APIKey == "" {
		return nil, errors.New("API key must not be empty")
	}
	return &Client{
		APIKey: APIKey,
		APIURL: "https://api.openweathermap.org",
		HTTPClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}, nil
}

func (c *Client) GetWeather(location string) (Conditions, error) {
	URL := fmt.Sprintf("%s/data/2.5/weather?q=%s&appid=%s", c.APIURL, url.QueryEscape(location), url.QueryEscape(c.APIKey))
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return Conditions{}, fmt.Errorf("error making API request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return Conditions{}, fmt.Errorf("unexpected response status: %v", resp.Status)
	}
	var data APIResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return Conditions{}, err
	}
	return Conditions{
		Summary:     data.Weather[0].Main,
		TemperatureKelvin: data.Main.Temp,
	}, nil
}

func ParseArgs(args []string) (fahrMode bool, location string) {
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.BoolVar(&fahrMode, "fahr", false, "Temperature in Fahrenheit")
	fs.Parse(args)
	location = fs.Arg(0)
	return fahrMode, location
}

func RunCLI() {
	fahrMode, location := ParseArgs(os.Args[1:])
	if location == "" {
		log.Fatal("Usage: weather LOCATION")
	}
	APIKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	client, err := NewClient(APIKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("getting weather for", location)
	conditions, err := client.GetWeather(location)
	if err != nil {
		log.Fatal(err)
	}
	if fahrMode {
		fmt.Println(conditions.StringFahrenheit())
	} else {
		fmt.Println(conditions.StringCelsius())
	}
}
