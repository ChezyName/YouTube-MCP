package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	YouTubeAPI    string //for youtube analytics
	ChannelHandle string
}

var cfg *Config

// Load the ENV
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg = &Config{
		YouTubeAPI:    os.Getenv("YouTubeAPI"),
		ChannelHandle: os.Getenv("ChannelHandle"),
	}
}

func GetConfig() *Config {
	return cfg
}
