/*
* Installs the program by using the GitHub releases to download the correct version based on your system,
* Sets up the config and all related files
* Then starts the key generation
* Finally returns the file location for the MCP executable
* Potentially generates the MCP code to allow copy & paste to AI config.json
 */

package main

import (
	"fmt"
	"os"

	"github.com/ChezyName/YouTube-MCP/config"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var headerView = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("62")).
	Padding(0, 1).
	Render("YouTube MCP Installation")

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		state:    []string{"Initilizing YouTube MCP", ""},
		spinner:  s,
		progress: progress.New(progress.WithDefaultGradient()),
	}
}

// Init runs first for the cmd
func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		checkConfig,
	)
}

// --- Update (The Logic) ---

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		//program is done from here
		if m.configStep == stateDone {
			switch msg.String() {
			case "enter", "ctrl+c", "q":
				return m, tea.Quit
			}
			return m, nil
		}

		// IF the user is currently typing in a text field (Step 2 or Step 3)
		switch m.configStep {
		case stateAPI, stateHandle, stateRequestHandleChange, stateReqDownload:
			if msg.String() == "enter" {
				return m.handleInputSubmission()
			}
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}

		if m.authStep == authStepClientID || m.authStep == authStepClientSecret {
			if msg.String() == "enter" {
				userInput := m.textInput.Value()

				if m.authStep == authStepClientID {
					m.tempClientID = userInput

					// Move to next input block: Client Secret
					m.authStep = authStepClientSecret
					m.state = append(m.state, "", "✓ Client ID stored.")
					m.state = append(m.state, "Please enter your Google OAuth Client Secret:")

					m.textInput.SetValue("")
					m.textInput.Placeholder = "GOCSPX-xxxxxxxxx"
					return m, nil
				}

				if m.authStep == authStepClientSecret {
					m.tempClientSecret = m.textInput.Value()
					m.authStep = authStepWaitingBrowser
					m.textInput.Blur()

					m.state[len(m.state)-1] = "✓ Client Secret stored."
					m.state = append(m.state, "(OAuth URL Will Be Copied - incase URL does not auto open)")
					m.state = append(m.state, "Launching authentication server...")

					// Fire the background callback server command asynchronously
					return m, startLocalOAuthServerCmd(m.tempClientID, m.tempClientSecret)
				}
			}
			var cmd tea.Cmd
			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case fatalError:
		if msg.err != nil {
			m.configStep = stateDone
			m.state = append(m.state, fmt.Sprintf("[FATAL ERROR]: %s", msg.err.Error()), "(press enter to quit)")
			return m, nil
		}
		return m, nil

	case configSetup:
		//start config setup
		return m.advanceSetupWizard()

	case loginFinishedMsg:
		// Auth finished in background. Move to the next step automatically.
		return m.advanceSetupWizard()

	case checkFinishedMsg:
		m.state = append(m.state,
			fmt.Sprintf("Config Fully Initialized @%s", config.GetConfig().ChannelHandle),
			fmt.Sprintf("Prepping Download for YouTube MCP (%s)", getOS()),
		)
		return m, nil
	case versionCheck:
		m.state = append(m.state, "") //create spacing regardless

		// Scenario C: No upstream version found, or it matches current version
		if msg.UpVersion == "" || msg.UpVersion == msg.CurrentVersion {
			m.state = append(m.state, "You are all set! YouTube-MCP is up to date.", "MCP Located at", getFileOut())
			m.configStep = stateDone
			return m, nil
		}

		// Scenario B: No local version found -> Auto-download
		if msg.CurrentVersion == "" || msg.CurrentVersion == "NO_FILE_FOUND" {
			m.configStep = stateDownload
			m.state = append(m.state, fmt.Sprintf("No local version found. Auto-downloading latest version (%s)...", msg.UpVersion))

			// Return a command to immediately start downloading in the background
			m.progressChan = make(chan float64, 100)
			return m, tea.Batch(
				downloadMCPCmd(m.progressChan),         // writes
				listeDownloadnProgress(m.progressChan), // reads
			)
		}

		// Scenario A: Both versions exist and don't match -> Prompt user
		if msg.CurrentVersion != "" && msg.UpVersion != "" {
			m.configStep = stateReqDownload
			msg := fmt.Sprintf("Update available! Current: %s -> Latest: %s", msg.CurrentVersion, msg.UpVersion)
			if m.state[len(m.state)-1] != msg {
				m.state = append(m.state, msg)
			}

			ti := textinput.New()
			ti.Placeholder = "y/n"
			ti.Focus()
			m.textInput = ti
			return m, nil
		}
	case downloadProgressMsg:
		m.downloadPct = float64(msg)
		progressCmd := m.progress.SetPercent(float64(msg * 100.0))
		return m, tea.Batch(
			progressCmd,
			listeDownloadnProgress(m.progressChan),
		)
	case progress.FrameMsg:
		newProgressModel, newCmd := m.progress.Update(msg)
		if pm, ok := newProgressModel.(progress.Model); ok {
			m.progress = pm
		}
		cmd = newCmd
		return m, cmd
	case downloadFinishedMsg:
		m.configStep = stateNone
		m.downloadPct = 1
		m.state = append(m.state, fmt.Sprintf("Downloaded to %s", getFileOut()))
		m.configStep = stateDone
		return m, nil
	}

	/*
		case checkFinishedMsg:
			m.configDir = msg.configPath
			m.state = stateDownloading
			return m, downloadFileCmd()

		case downloadFinishedMsg:
			m.state = stateConfiguring
			return m, nil

		case progressMsg:
			m.downloadPct = float64(msg)
			if m.downloadPct >= 1.0 {
				m.state = stateConfiguring
			}
			return m, nil
	*/

	return m, cmd
}

// --- View (The UI) ---

func (m model) View() string {
	var bodyView string

	// Build rolling scroll history list log block layout
	var logHistory string
	for _, logLine := range m.state {
		logHistory += fmt.Sprintf("   %s\n", logLine)
	}

	switch m.configStep {
	case stateDownload:
		bodyView = fmt.Sprintf("%s\n   Downloading YouTube MCP (%s):\n   %s\n", getOS(), logHistory, m.progress.View())
	case stateAPI, stateHandle, stateRequestHandleChange, stateReqDownload:
		bodyView = fmt.Sprintf("%s\n%s\n\n(Press Enter to confirm)", logHistory, m.textInput.View())
	default:
		switch m.authStep {
		case authStepClientID, authStepClientSecret:
			// If actively prompt typing: show logs along with live flashing keyboard box
			bodyView = fmt.Sprintf("%s\n%s\n\n(Press Enter to confirm)", logHistory, m.textInput.View())
		case authStepWaitingBrowser:
			// If waiting for web redirect click: show progress spinner inline on last log trace
			bodyView = fmt.Sprintf("%s\n%s Waiting for browser auth confirmation link callback...", logHistory, m.spinner.View())
		default:
			bodyView = logHistory
		}
	}

	return fmt.Sprintf("%s\n%s", headerView, bodyView)
}

func main() {
	m := initialModel()
	p := tea.NewProgram(m)

	if _, err := p.Run(); err != nil {
		fmt.Printf("An error has occured: %v", err)
		os.Exit(1)
	}
}

func (m model) handleInputSubmission() (model, tea.Cmd) {
	userInput := m.textInput.Value()
	cfg := config.GetConfig()

	switch m.configStep {
	case stateAPI:
		cfg.YouTubeAPI = userInput
	case stateHandle:
		cfg.ChannelHandle = userInput
	case stateRequestHandleChange:
		if userInput == "y" || userInput == "Y" {
			// Drop into edit mode
			m.configStep = stateHandle
			m.state = append(m.state, "Enter your new Channel Handle (without the @):")

			ti := textinput.New()
			if SuggestedChannelHandle != "" {
				ti.SetValue(SuggestedChannelHandle)
			}
			ti.Placeholder = "YourChannel"
			ti.Focus()
			m.textInput = ti
			return m, nil
		} else {
			// Skip — treat as already set, advance
			m.configStep = stateNone
			m.textInput.Blur()
			return m.advanceSetupWizard()
		}
	case stateReqDownload:
		if userInput == "y" || userInput == "Y" {
			m.configStep = stateDownload
			m.progressChan = make(chan float64, 100)
			return m, tea.Batch(
				downloadMCPCmd(m.progressChan),         // writes
				listeDownloadnProgress(m.progressChan), // reads
			)
		} else {
			// User skipped download, there good to go
			m.configStep = stateNone
			m.textInput.Blur()
			m.state = append(m.state, "You are all set!")
			return m, tea.Quit
		}
	}

	saveConfig(*cfg)
	m.textInput.Blur()
	return m.advanceSetupWizard()
}

func listeDownloadnProgress(progressChan chan float64) tea.Cmd {
	return func() tea.Msg {
		percent, ok := <-progressChan
		if !ok {
			return nil
		}

		return downloadProgressMsg(percent)
	}
}
