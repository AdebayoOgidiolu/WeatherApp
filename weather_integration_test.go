// +build integration

package weather_test

import (
	"os"
	"testing"
	"weather"
)

func TestConditionsIntegration(t *testing.T) {
	APIKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	if APIKey == "" {
		t.Fatal("OPENWEATHERMAP_API_KEY must be set to run this test")
	}
	client, err := weather.NewClient(APIKey)
	if err != nil {
		t.Fatal(err)
	}
	conditions, err := client.Current("London")
	if err != nil {
		t.Fatal(err)
	}
	if conditions.Summary == "" {
		t.Error("want non-empty summary")
	}
}