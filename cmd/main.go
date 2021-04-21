package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"weather"

	"github.com/enescakir/emoji"
)

func main() {

	APIKey := os.Getenv("OPENWEATHERMAP_API_KEY")
	fahr := flag.Bool("fahr", false, "Temperature in Fahrenheit")
	flag.Parse()

	location := os.Args[1]
	client, err := weather.NewClient(APIKey)
	if err != nil {
		log.Fatal(err)
	}
	conditions, err := client.GetWeather(location, fahr)
	fmt.Println(emoji.Sun, " ", conditions.Summary, conditions.Temperature, "C")
}
