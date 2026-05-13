package config

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

func EnsureRefreshToken() {
	if os.Getenv("YOUTUBE_REFRESH_TOKEN") != "" {
		return
	}

	conf := GetOAuthConfig()
	authURL := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	fmt.Println("========================================")
	fmt.Println("No refresh token found!")
	fmt.Println("Visit this URL to authorize:")
	fmt.Println()
	fmt.Println(authURL)
	fmt.Println()
	fmt.Println("Waiting for authorization...")
	fmt.Println("This will automatically update ENV if possible.")
	fmt.Println("========================================")

	done := make(chan struct{})

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			log.Fatal("Token exchange failed:", err)
		}
		saveRefreshToken(token.RefreshToken)
		os.Setenv("YOUTUBE_REFRESH_TOKEN", token.RefreshToken)

		w.Write([]byte("<h2>Authorization successful! You can close this tab and return to your server.</h2>"))
		close(done)
	})

	srv := &http.Server{Addr: ":9999", Handler: mux}
	go srv.ListenAndServe()

	<-done
	srv.Shutdown(context.Background())

	fmt.Println("Authorization complete! Refresh token saved to .env")
}

func saveRefreshToken(token string) {
	envMap, err := godotenv.Read(".env")
	if err != nil {
		envMap = make(map[string]string)
	}

	envMap["YOUTUBE_REFRESH_TOKEN"] = token

	if err := godotenv.Write(envMap, ".env"); err != nil {
		fmt.Printf("\nCould not write to .env — add this manually:\nYOUTUBE_REFRESH_TOKEN=%s\n", token)
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
