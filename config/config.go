package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	YouTubeAPI          string //for youtube analytics
	YouTubeRefreshToken string //for youtube analytics
	ChannelHandle       string
}

var cfg *Config

// Load the ENV
func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system variables")
	}

	cfg = &Config{
		YouTubeAPI:          os.Getenv("YouTubeAPI"),
		ChannelHandle:       os.Getenv("ChannelHandle"),
		YouTubeRefreshToken: os.Getenv("YOUTUBE_REFRESH_TOKEN"),
	}
}

func GetConfig() *Config {
	return cfg
}
