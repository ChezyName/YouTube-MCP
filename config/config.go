package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	YouTubeAPI            string `json:"YOUTUBE_API"`
	YouTubeRefreshToken   string `json:"YOUTUBE_REFRESH_TOKEN"`
	ChannelHandle         string `json:"ChannelHandle"`
	YOUTUBE_CLIENT_ID     string `json:"YOUTUBE_CLIENT_ID"`
	YOUTUBE_CLIENT_SECRET string `json:"YOUTUBE_CLIENT_SECRET"`
}

var cfg *Config

// Load from %appdata%/YouTube-MCP/config.json
// or create it and tell user to set the params
func LoadConfig() {
	appData, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR Trying to look for appdata path:%s", err.Error())
		os.Exit(1)
	}

	configDir := filepath.Join(appData, "/YouTube-MCP")
	configFile := filepath.Join(configDir, "config.json")

	//make the folder if does not exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, 0755)
	}

	//does not exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultCfg := Config{}

		bytes, _ := json.MarshalIndent(defaultCfg, "", "  ")

		err := os.WriteFile(configFile, bytes, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Could not create config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Fprintf(os.Stderr, "\n[CONFIG REQUIRED]\n")
		fmt.Fprintf(os.Stderr, "A default config file has been created at:\n%s\n", configFile)
		fmt.Fprintf(os.Stderr, "Please open this file, add your API keys, and restart the server.\n\n")
		fmt.Fprintf(os.Stderr, "Use the token-getter to get a token for the 'YOUTUBE_REFRESH_TOKEN' after updating your keys..\n\n")
		os.Exit(0)
	}

	file, err := os.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not read config file: %v\n", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(file, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Invalid JSON in config file: %v\n", err)
		os.Exit(1)
	}
}

func GetConfig() *Config {
	return cfg
}
