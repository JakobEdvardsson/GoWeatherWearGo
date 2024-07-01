package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/JakobEdvardsson/GoWeatherWearGo/api"
	"github.com/JakobEdvardsson/GoWeatherWearGo/storage"
	"github.com/JakobEdvardsson/GoWeatherWearGo/util"
)

func init() {
	weatherApiKey := os.Getenv("API-KEY-WEATHERAPI")
	if weatherApiKey == "" {
		err := util.LoadEnvFile(".env")
		if err != nil {
			log.Fatal("No env file or env vars provided!")
		}
	}
}

func main() {
	listenPort := flag.String("p", "8080", "The listening port. (Default 8080)")
	flag.Parse()

	fmt.Printf("Running on port: %v\n", *listenPort)

	storage := storage.NewPostgresStorage()

	weatherApiKey := os.Getenv("API-KEY-WEATHERAPI")
	if weatherApiKey == "" {
		err := util.LoadEnvFile(".env")
		if err != nil {
			log.Fatal("No env file or env vars provided!")
		}
	}

	server := api.NewServer(*listenPort, storage, os.Getenv("API-KEY-WEATHERAPI"))
	log.Fatal(server.Start())
}
