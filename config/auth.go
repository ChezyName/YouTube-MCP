package config

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func GetOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     GetConfig().YOUTUBE_CLIENT_ID,
		ClientSecret: GetConfig().YOUTUBE_CLIENT_SECRET,
		RedirectURL:  "http://localhost:9999/callback",
		Scopes: []string{
			youtube.YoutubeReadonlyScope,
			"https://www.googleapis.com/auth/yt-analytics.readonly",
		},
		Endpoint: google.Endpoint,
	}
}

func GetOAuthClient() (*http.Client, error) {
	cfg := GetOAuthConfig()
	if cfg == nil {
		return nil, fmt.Errorf("Invalid config, Issue with config.json or project")
	}

	if GetConfig().YouTubeRefreshToken == "" {
		return nil, fmt.Errorf("Invalid YouTubeRefreshToken, Please edit the config.json")
	}

	token := &oauth2.Token{
		RefreshToken: GetConfig().YouTubeRefreshToken,
	}
	tokenSource := cfg.TokenSource(context.Background(), token)
	return oauth2.NewClient(context.Background(), tokenSource), nil
}
