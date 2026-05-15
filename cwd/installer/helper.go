package main

import (
	"errors"
	"os"
	"path/filepath"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}

	return false
}

// configDir & configFile Locations
func getConfigDir() (string, string) {
	appData, _ := os.UserConfigDir()
	configDir := filepath.Join(appData, "YouTube-MCP")
	configFile := filepath.Join(configDir, "config.json")

	return configDir, configFile
}
