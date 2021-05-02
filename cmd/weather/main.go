package main

import (
	"log"
	"os"

	"github.com/aculclasure/weather"
)

func main() {
	if err := weather.CurrentWeatherCLI(os.Args); err != nil {
		log.Fatal(err)
	}
}
