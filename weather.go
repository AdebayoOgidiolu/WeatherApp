package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
	Temperature float64
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
		Temperature: data.Main.Temp - 273.15,
	}, nil
}
