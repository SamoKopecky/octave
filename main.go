package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/lukasl-dev/octave/config"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Failed to get .env file: %s", err)
	}
	app := newApp(config.Config{
		Token: os.Getenv("OCTAVE_TOKEN"),
		Lavalink: config.Lavalink{
			Host:       getEnv("LAVALINK_HOST", "127.0.0.1:2334"),
			Passphrase: getEnv("LAVALINK_PASSPHRASE", "alpharius"),
		},
	})

	if err := app.run(); err != nil {
		log.Fatalf("Failed to run app: %s\n", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, os.Kill)
	<-signals
}

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

// TODO: options for easier running
