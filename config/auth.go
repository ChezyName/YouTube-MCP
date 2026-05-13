package config

import (
	"context"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func GetOAuthConfig() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("YOUTUBE_CLIENT_ID"),
		ClientSecret: os.Getenv("YOUTUBE_CLIENT_SECRET"),
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
	token := &oauth2.Token{
		RefreshToken: GetConfig().YouTubeRefreshToken,
	}
	tokenSource := cfg.TokenSource(context.Background(), token)
	return oauth2.NewClient(context.Background(), tokenSource), nil
}
