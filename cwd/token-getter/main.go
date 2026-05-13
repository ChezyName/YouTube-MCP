package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.design/x/clipboard"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func main() {
	fmt.Println("========================================")
	fmt.Println("  YouTube MCP - Token Getter")
	fmt.Println("========================================")
	fmt.Println()

	// Get credentials from user
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("This following should not be shared - especially the refresh token.\n")
	fmt.Print("Client ID: ")
	clientID, _ := reader.ReadString('\n')
	clientID = strings.TrimSpace(clientID)

	fmt.Print("Client Secret: ")
	clientSecret, _ := reader.ReadString('\n')
	clientSecret = strings.TrimSpace(clientSecret)

	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:9999/callback",
		Scopes: []string{
			youtube.YoutubeReadonlyScope,
			"https://www.googleapis.com/auth/yt-analytics.readonly",
		},
		Endpoint: google.Endpoint,
	}

	authURL := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)

	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("Visit this URL to authorize:")
	fmt.Println()
	fmt.Println(authURL)
	fmt.Println()
	fmt.Println("Waiting for authorization...")
	fmt.Println("Will attempt to save token to clipboard.")
	fmt.Println("========================================")

	done := make(chan struct{})
	doneToken := make(chan struct{})

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

		saveToken(token.RefreshToken, doneToken)

		w.Write([]byte("Authorization successful! You can close this tab and return to your terminal."))
		close(done)
	})

	srv := &http.Server{Addr: ":9999", Handler: mux}
	go srv.ListenAndServe()

	<-done
	srv.Shutdown(context.Background())

	<-doneToken
	fmt.Println("Authorization complete!")
}

func saveToken(token string, done chan struct{}) {
	// Copy to clipboard
	if err := clipboard.Init(); err == nil {
		clipboard.Write(clipboard.FmtText, []byte(token))
		fmt.Println("Copied to clipboard!")
	} else {
		fmt.Println("Could not copy to clipboard")
	}

	// Ask in a goroutine so it doesn't block the HTTP response
	go func() {
		fmt.Print("\nWould you like to print the refresh token? (note: this is secret and should not be shared) [y/n]: ")
		var answer string
		fmt.Scanln(&answer)
		if answer == "y" || answer == "Y" {
			fmt.Println(token)
		} else {
			fmt.Println("Token not printed.")
		}

		close(done)
	}()
}
