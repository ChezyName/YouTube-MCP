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
		state:       []string{"Initilizing YouTube MCP"},
		spinner:     s,
		progress:    progress.New(progress.WithDefaultGradient()),
		downloadPct: 0,
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

		// IF the user is currently typing in a text field (Step 2 or Step 3)
		if m.configStep == stateAPI || m.configStep == stateHandle {
			if msg.String() == "enter" {
				return m.handleInputSubmission()
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
		m.state = append(m.state, msg.err.Error())
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
	if m.configStep == stateAPI || m.configStep == stateHandle {
		var logHistory string
		for _, logLine := range m.state {
			logHistory += fmt.Sprintf("  %s\n", logLine) //indent for all to fix spinner pos issue
		}

		bodyView = fmt.Sprintf(
			"%s\n%s\n\n(Press Enter to confirm)",
			logHistory,
			m.textInput.View(),
		)
	} else {
		var logHistory string
		totalLogs := len(m.state)

		for i, logLine := range m.state {
			//current active log
			if i == totalLogs-1 {
				logHistory += fmt.Sprintf("%s %s\n", m.spinner.View(), logLine)
			} else {
				logHistory += fmt.Sprintf("✓ %s\n", logLine)
			}
		}
		bodyView = logHistory
	}

	// Stack the persistent header cleanly above our custom body view
	return fmt.Sprintf("%s\n%s", headerView, bodyView)
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error running program:", err)
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
	}

	saveConfig(*cfg)
	m.textInput.Blur()
	return m.advanceSetupWizard()
}
