package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/ChezyName/YouTube-MCP/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Checks the Config
func checkConfig() tea.Msg {
	// Simulate a slight delay for "looking good"
	configDir, configFile := getConfigDir()

	// Create the directory if it doesn't exist
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return fatalError{err: err}
		}
	}

	if !fileExists(configFile) {
		_, err := os.Create(configFile)
		if err != nil {
			return fatalError{err: err}
		}

		//save empty
		saveConfig(config.Config{})
		return configSetup{}
	}

	//Make sure the JSON is a valid structure or it will fail
	file, err := os.ReadFile(configFile)
	if err != nil {
		saveConfig(config.Config{})
	}

	var cfg config.Config
	if err := json.Unmarshal(file, &cfg); err != nil {
		saveConfig(config.Config{})
	}

	//Load Config first
	config.LoadConfig()

	//Start setup to check auth anyways
	return configSetup{}
}

func loadConfig() *config.Config {
	config.LoadConfig()
	return config.GetConfig()
}

func saveConfig(conf config.Config) {
	_, file := getConfigDir()
	fileData, _ := json.MarshalIndent(conf, "", "  ")
	_ = os.WriteFile(file, fileData, 0644)
}

var lastCheckAuthResult *bool = nil
var lastCheckAPIResult *bool = nil

var passedAuth = false
var passedAPI = false
var passedHandle = false

var SuggestedChannelHandle = ""
var suggestedChannelHandleOnce sync.Once

func checkAuth(cfg *config.Config) bool {
	if cfg.YouTubeRefreshToken == "" || (cfg.YOUTUBE_CLIENT_ID == "" || cfg.YOUTUBE_CLIENT_SECRET == "") {
		return false
	}

	if lastCheckAuthResult != nil && *lastCheckAuthResult == true {
		return *lastCheckAuthResult
	}

	//test to see if end-point can be reached
	client, err := config.GetOAuthClient()
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// We use an authenticated request to the YouTube metadata validation API.
	// Filtering by mine=true uses almost zero quota but requires operational OAuth credentials.
	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://www.googleapis.com/youtube/v3/channels?part=snippet,brandingSettings&mine=true", nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	isWorking := resp.StatusCode == http.StatusOK

	// Parse the handle out of the response
	if isWorking {
		var result youtube.ChannelListResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
			if len(result.Items) > 0 {
				suggestedChannelHandleOnce.Do(func() {
					SuggestedChannelHandle = result.Items[0].Snippet.Title
				})
			}
		}
	}

	lastCheckAuthResult = &isWorking
	return isWorking
}

func checkAPI(m *model, cfg *config.Config) bool {
	if cfg == nil {
		return false
	}

	if cfg.YouTubeAPI == "" {
		return false
	}

	if lastCheckAPIResult != nil && *lastCheckAPIResult == true {
		return *lastCheckAPIResult
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	svc, err := youtube.NewService(ctx, option.WithAPIKey(cfg.YouTubeAPI))
	if err != nil {
		return false
	}

	_, err = svc.VideoCategories.List([]string{"snippet"}).RegionCode("US").Do()

	isWorking := err == nil
	lastCheckAPIResult = &isWorking
	return isWorking
}

func (m model) advanceSetupWizard() (model, tea.Cmd) {

	// STEP 1: AUTHENTICATION
	cfg := loadConfig()
	if !checkAuth(cfg) {
		if m.authStep == authStepNone {
			// Trigger Step 1A: Prompt for Client ID
			m.authStep = authStepClientID
			if m.state[len(m.state)-1] != instructions {
				m.state = append(m.state, instructions)
			}

			ti := textinput.New()
			ti.Placeholder = "123456-abc.apps.googleusercontent.com"
			ti.Focus()
			m.textInput = ti
			return m, nil
		}
	} else if !passedAuth {
		passedAuth = true
		m.state = append(m.state, "OAuth Succsesfully Passed", "")
		m.authStep = authStepNone
	}

	// STEP 2: PUBLIC API KEY
	cfg = loadConfig()
	if !checkAPI(&m, cfg) {
		m.configStep = stateAPI
		var msg = "Step 2/3: Enter your YouTube Public API Key:"
		if m.state[len(m.state)-1] != msg {
			m.state = append(m.state, msg)
		}

		ti := textinput.New()
		ti.Placeholder = "AIzaSy..."
		ti.Focus()
		m.textInput = ti
		return m, nil // Stop background work; wait for input
	} else if !passedAPI {
		passedAPI = true
		m.state = append(m.state, "API Check Succsesfully Passed", "")
	}

	// STEP 3: CHANNEL HANDLE
	cfg = loadConfig()
	if cfg.ChannelHandle == "" {
		passedHandle = true //cannot edit if alr editing
		m.configStep = stateHandle
		var msg = "Step 3/3: Enter your YouTube Channel Handle (without the @):"
		if m.state[len(m.state)-1] != msg {
			m.state = append(m.state, msg)
		}

		ti := textinput.New()
		if SuggestedChannelHandle != "" {
			ti.SetValue(SuggestedChannelHandle)
		}
		ti.Placeholder = "YourChannel"
		ti.Focus()
		m.textInput = ti
		return m, nil
	} else if m.configStep != stateRequestHandleChange && !passedHandle {
		passedHandle = true
		m.configStep = stateRequestHandleChange
		msg := fmt.Sprintf("Step 3/3: Channel Handle is already set to @%s. Change it? [y/n]:", cfg.ChannelHandle)
		if m.state[len(m.state)-1] != msg {
			m.state = append(m.state, msg)
		}

		ti := textinput.New()
		ti.Placeholder = "y/n"
		ti.Focus()
		m.textInput = ti
		return m, nil
	}

	// STEP 4: DONE
	m.configStep = stateNone
	var msg = "All configuration verifications passed successfully!"
	if m.state[len(m.state)-1] != msg {
		m.state = append(m.state, msg)
	}

	// return command for
	return m, downloadMCPCmd()
}
