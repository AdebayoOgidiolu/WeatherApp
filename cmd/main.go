package main

import (
	"fmt"
	"log"
	"os"
	"weather"
)

func main() {
	APIKey := os.Getenv("OPENWEATHERMAP_API_KEY")

	location := os.Args[1]
	client, err := weather.NewClient(APIKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client.GetWeather(location))
}
