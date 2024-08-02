package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/JakobEdvardsson/GoWeatherWearGo/api"
	"github.com/JakobEdvardsson/GoWeatherWearGo/storage"
	"github.com/JakobEdvardsson/GoWeatherWearGo/util"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type EnvVars struct {
	weatherApiKey              string
	AUTH_SPOTIFY_CLIENT_ID     string
	AUTH_SPOTIFY_CLIENT_SECRET string
	dbHost                     string
	dbPort                     int
	dbUsername                 string
	dbPassword                 string
	dbDbname                   string
}

var envVars *EnvVars

func init() {
	if util.CheckEnvFileExists(".env") {
		err := util.LoadEnvFile(".env")
		if err != nil {
			log.Fatal("main.go: incomplete .env file")
		}
	}

	envVars.weatherApiKey = os.Getenv("API_KEY_WEATHERAPI")
	envVars.AUTH_SPOTIFY_CLIENT_ID = os.Getenv("AUTH_SPOTIFY_CLIENT_ID")
	envVars.AUTH_SPOTIFY_CLIENT_SECRET = os.Getenv("AUTH_SPOTIFY_CLIENT_SECRET")
	envVars.dbHost = os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")

	var err error
	envVars.dbPort, err = strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatal("Could not convert DB_PORT to int")
	}

	envVars.dbUsername = os.Getenv("DB_USER")
	envVars.dbPassword = os.Getenv("DB_PW")
	envVars.dbDbname = os.Getenv("DB_NAME")

	if envVars.weatherApiKey == "" ||
		envVars.AUTH_SPOTIFY_CLIENT_ID == "" ||
		envVars.AUTH_SPOTIFY_CLIENT_SECRET == "" ||
		envVars.dbHost == "" ||
		dbPortStr == "" ||
		envVars.dbUsername == "" ||
		envVars.dbPassword == "" ||
		envVars.dbDbname == "" {

		log.Fatal("main.go: missing required environment variables")
	}
}

func main() {
	listenPort := flag.String("p", "8080", "The listening port. (Default 8080)")
	flag.Parse()

	fmt.Printf("Running on port: %v\n", *listenPort)
	storage := storage.NewPostgresStorage(envVars.dbHost, envVars.dbUsername, envVars.dbPassword, envVars.dbDbname, envVars.dbPort)
	defer storage.DB.Close()

	weatherApiKey := os.Getenv("API_KEY_WEATHERAPI")
	if weatherApiKey == "" {
		err := util.LoadEnvFile(".env")
		if err != nil {
			log.Fatal("main.go: No env file or env vars provided!")
		}
	}
	ExampleClient()

	server := api.NewServer(*listenPort, storage, os.Getenv("API_KEY_WEATHERAPI"))
	log.Fatal(server.Start())
}

var ctx = context.Background()

func ExampleClient() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
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
