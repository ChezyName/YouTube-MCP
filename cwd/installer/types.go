package main

import (
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

// Messages define state changes
type checkFinishedMsg struct{}
type fatalError struct{ err error }
type configSetup struct{}
type loginFinishedMsg struct{} //for Auth

type globalState int

const (
	stateNone globalState = iota
	stateAuth
	stateHandle
	stateRequestHandleChange
	stateReqDownload
	stateDownload
	stateAPI
	stateDone //program will close after this - no return
)

type model struct {
	state      []string
	configStep globalState
	progress   progress.Model
	spinner    spinner.Model
	textInput  textinput.Model

	authStep         authSubStep
	tempClientID     string
	tempClientSecret string

	//for downloads
	progressChan chan float64
	downloadPct  float64
}
