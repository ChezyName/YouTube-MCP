package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var repoURL = "https://api.github.com/repos/ChezyName/YouTube-MCP/releases/latest"

type Release struct {
	Version   string  `json:"tag_name"`
	Assets    []asset `json:"assets"`
	AssetsMap map[string]string
}

type asset struct {
	Name     string `json:"name"`
	Download string `json:"browser_download_url"`
}

var currentRelease *Release

func getLatestDownload() (*Release, error) {
	if currentRelease != nil {
		return currentRelease, nil
	}

	// Set up a client with a timeout so it doesn't hang indefinitely
	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Get(repoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned Status: %s", http.StatusText(resp.StatusCode))
	}

	// Decode the JSON response directly into a new Release struct
	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("Could not decode JSON Resposne body")
	}

	//make the assets map by connecting asset name to the download url
	var AssetsMap = make(map[string]string, len(release.Assets))
	for _, asset := range release.Assets {
		AssetsMap[asset.Name] = asset.Download
	}

	release.AssetsMap = AssetsMap

	// Cache the result for future calls
	currentRelease = &release
	return currentRelease, nil
}
