package main

import (
	"fmt"
	"log"
	"weather"
)

func main() {
	client, err := weather.NewClient(APIKey)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(client.Conditions(location))
}

// type APIResponse struct{
// 	Weather []struct{
// 		Main string
// 	}
// 	Main struct{
// 		Temp float64
// 	}
// }

// type Conditions struct{
// 	Summary string
// 	Temperature float64
// }
//

// func main() {
// 	resp, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q=London&appid=3b814c61996538f2e8a2b921e23bbb0a")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer resp.Body.Close()
// 	if resp.StatusCode != http.StatusOK {
// 		fmt.Fprintf(os.Stderr, "unexpected response status: %v", resp.Status)
// 		os.Exit(1)
// 	}
// 	var data APIResponse
// 	err = json.NewDecoder(resp.Body).Decode(&data)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	conditions := Conditions{
// 		Summary: data.Weather[0].Main,
// 		Temperature: data.Main.Temp-273.15,
// 	}
// 	fmt.Printf("%+v\n", conditions)
// }