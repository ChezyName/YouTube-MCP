package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

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
	openBrowser(authURL)

	<-done
	srv.Shutdown(context.Background())

	<-doneToken
	fmt.Println("Authorization complete!")
}

func copyToClipboard(token string) bool {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("clip")
	case "darwin":
		cmd = exec.Command("pbcopy")
	default: // linux
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd = exec.Command("xclip", "-selection", "clipboard")
		} else if _, err := exec.LookPath("wl-clipboard"); err == nil {
			cmd = exec.Command("wl-copy")
		} else {
			return false
		}
	}
	cmd.Stdin = strings.NewReader(token)
	return cmd.Run() == nil
}

func saveToken(token string, done chan struct{}) {
	// Copy to clipboard
	if copyToClipboard(token) {
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

// OS-independent browser utility helper
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		// Try common browser paths directly instead of shell
		browsers := []string{
			`C:\Program Files\Google\Chrome\Application\chrome.exe`,
			`C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`,
			os.ExpandEnv(`%LOCALAPPDATA%\Google\Chrome\Application\chrome.exe`),
			`C:\Program Files\Mozilla Firefox\firefox.exe`,
			`C:\Program Files\Microsoft\Edge\Application\msedge.exe`,
		}
		for _, b := range browsers {
			if _, err := os.Stat(b); err == nil {
				return exec.Command(b, url).Start()
			}
		}
		// Fallback to shell
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	return exec.Command(cmd, args...).Start()
}
