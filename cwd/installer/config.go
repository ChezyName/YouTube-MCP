package main

import (
	"encoding/json"
	"os"

	"github.com/ChezyName/YouTube-MCP/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
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

		return configSetup{}
	}

	//Load Config first
	config.LoadConfig()

	//Better to check before fatal, but should never be invalid here
	if config.GetConfig() == nil {
		return configSetup{}
	}

	//Check the configs
	//Need to AUTH
	if config.GetConfig().YouTubeRefreshToken == "" ||
		(config.GetConfig().YOUTUBE_CLIENT_ID == "" || config.GetConfig().YOUTUBE_CLIENT_SECRET == "") {
		//Load Secrets & Auth
		return configSetup{}
	}

	//Need User Handle - Can Auto Get from Channel Analytics
	if config.GetConfig().ChannelHandle == "" {
		//Load Channel Handle
		return configSetup{}
	}

	//Need YouTube Public API Key
	if config.GetConfig().YouTubeAPI == "" {
		//Load Channel Handle
		return configSetup{}
	}

	return checkFinishedMsg{configPath: configDir}
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

func (m model) advanceSetupWizard() (model, tea.Cmd) {

	// STEP 1: AUTHENTICATION
	cfg := loadConfig()
	if cfg.YouTubeRefreshToken == "" && (cfg.YOUTUBE_CLIENT_ID == "" || cfg.YOUTUBE_CLIENT_SECRET == "") {
		m.configStep = stateAuth
		m.state = append(m.state, "Step 1/3: Launching browser for YouTube Authentication...")

		// Return your background auth function wrapped as a command
		return m, func() tea.Msg {
			return loginClientAuth()
		}
	}

	// STEP 2: PUBLIC API KEY
	cfg = loadConfig()
	if cfg.YouTubeAPI == "" {
		m.configStep = stateAPI
		m.state = append(m.state, "Step 2/3: Enter your YouTube Public API Key:")

		ti := textinput.New()
		ti.Placeholder = "AIzaSy..."
		ti.Focus()
		m.textInput = ti
		return m, nil // Stop background work; wait for input
	}

	// STEP 3: CHANNEL HANDLE
	cfg = loadConfig()
	if cfg.ChannelHandle == "" {
		m.configStep = stateHandle
		m.state = append(m.state, "Step 3/3: Enter your YouTube Channel Handle (without the @):")

		ti := textinput.New()
		ti.Placeholder = "YourChannel"
		ti.Focus()
		m.textInput = ti
		return m, nil
	}

	// STEP 4: DONE
	m.configStep = stateNone
	m.state = append(m.state, "All configuration verifications passed successfully!")

	// return command for
	return m, nil
}
