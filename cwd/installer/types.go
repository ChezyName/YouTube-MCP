package main

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

// Messages define state changes
type checkFinishedMsg struct{ configPath string }
type fatalError struct{ err error }
type configSetup struct{}
type loginFinishedMsg struct{} //for Auth

type configState int

const (
	stateNone configState = iota
	stateAuth
	stateHandle
	stateAPI
)

type model struct {
	state       []string
	configStep  configState
	progress    progress.Model
	spinner     spinner.Model
	textInput   textinput.Model
	downloadPct float64
}
