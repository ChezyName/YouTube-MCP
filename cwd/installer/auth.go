package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/ChezyName/YouTube-MCP/config"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

func makeLink(text, url string) string {
	// Standard ANSI escape code to create true clickable links in modern terminals
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}

var oAuth = makeLink("OAuth Screen", "https://console.cloud.google.com/apis/credentials/consent")
var youtubeAPI = makeLink("YouTube API", "https://console.cloud.google.com/apis/library/youtube.googleapis.com")
var analyticsAPI = makeLink("YouTube Analytics API", "https://console.cloud.google.com/apis/library/youtubeanalytics.googleapis.com")
var instructions = fmt.Sprintf("You need to create an %s with permissions to the %s and %s.", oAuth, youtubeAPI, analyticsAPI)

type authSubStep int

const (
	authStepNone           authSubStep = iota
	authStepClientID                   // User is typing Client ID
	authStepClientSecret               // User is typing Client Secret
	authStepWaitingBrowser             // Waiting for browser click callback loop
)

func startLocalOAuthServerCmd(clientID, clientSecret string) tea.Cmd {
	return func() tea.Msg {
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
		copyToClipboard(authURL)

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

			cfg := config.GetConfig()
			if cfg == nil {
				cfg = &config.Config{}
			}

			cfg.YOUTUBE_CLIENT_ID = clientID
			cfg.YOUTUBE_CLIENT_SECRET = clientSecret
			cfg.YouTubeRefreshToken = token.RefreshToken
			saveConfig(*cfg)

			w.Write([]byte("Authorization successful! You can safely close this browser tab and return to your terminal program."))

			// Notify main routine thread execution loop to unwind server
			close(done)
		})

		srv := &http.Server{Addr: ":9999", Handler: mux}
		go srv.ListenAndServe()

		// Wait for OS to open port
		for {
			conn, err := net.Dial("tcp", "localhost:9999")
			if err == nil {
				conn.Close()
				break
			}
			time.Sleep(50 * time.Millisecond)
		}

		_ = openBrowser(authURL)

		// Block and sleep here inside the command until browser interaction resolves
		<-done
		_ = srv.Shutdown(context.Background())

		// Dispatch message back to main thread bubble tea loop safely
		return loginFinishedMsg{}
	}
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
