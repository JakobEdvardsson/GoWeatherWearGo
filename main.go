package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/JakobEdvardsson/GoWeatherWearGo/api"
	"github.com/JakobEdvardsson/GoWeatherWearGo/storage"
	"github.com/JakobEdvardsson/GoWeatherWearGo/util"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// TODO: Replace hardcoded values with env vars in application
func init() {
	weatherApiKey := os.Getenv("API_KEY_WEATHERAPI")
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
	defer storage.DB.Close()

	weatherApiKey := os.Getenv("API_KEY_WEATHERAPI")
	if weatherApiKey == "" {
		err := util.LoadEnvFile(".env")
		if err != nil {
			log.Fatal("No env file or env vars provided!")
		}
	}
	ExampleClient()

	server := api.NewServer(*listenPort, storage, os.Getenv("API_KEY_WEATHERAPI"))
	log.Fatal(server.Start())
}

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "host.docker.internal:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	err := rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}
