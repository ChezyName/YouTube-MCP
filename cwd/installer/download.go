package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// File should be config dir/YouTubeMCP
func getFileOut() string {
	appDir, _ := getConfigDir()
	var binaryFile = "YouTube-MCP"

	//should just be the binary except windows req .exe
	os := runtime.GOOS
	switch os {
	case "windows":
		binaryFile = "YouTube-MCP.exe"
	}

	return filepath.Join(appDir, binaryFile)
}

func getOS() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch os {
	case "windows":
		return "windows"

	case "linux":
		if arch == "arm64" {
			return "linux-arm64"
		}
		return "linux-amd64"

	case "darwin": // macOS
		if arch == "arm64" {
			return "darwin-arm64" // Apple Silicon (M1/M2/M3/M4)
		}
		return "darwin-amd64" // Intel Macs

	default:
		//guess choice
		return "unknown"
	}
}

func getDownloadFile() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch os {
	case "windows":
		// Windows only needs amd64 based on build script
		return "youtube-mcp-windows.exe"

	case "linux":
		if arch == "arm64" {
			return "youtube-mcp-linux-arm64"
		}
		// Default fallback for Linux is amd64
		return "youtube-mcp-linux-amd64"

	case "darwin": // macOS
		if arch == "arm64" {
			return "youtube-mcp-darwin-arm64" // Apple Silicon (M1/M2/M3/M4)
		}
		return "youtube-mcp-darwin-amd64" // Intel Macs

	default:
		// Fallback safe choice if something wild happens
		return "youtube-mcp-linux-amd64"
	}
}

// Run file with -v if able, then take the output
var currentVersion = ""

func getCurrentVersion() (string, error) {
	if currentVersion != "" {
		return currentVersion, nil
	}

	binaryPath := getFileOut()
	if !fileExists(binaryPath) {
		return "NO_FILE_FOUND", fmt.Errorf("Executable was not found for %s", getOS())
	}

	cmd := exec.Command(binaryPath, "-v")

	var out bytes.Buffer
	cmd.Stdout = &out

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return "UNKNOWN", fmt.Errorf("failed to run command: %v, stderr: %s", err, stderr.String())
	}

	currentVersion = strings.TrimSpace(out.String())
	return currentVersion, nil
}

type versionCheck struct {
	CurrentVersion string
	UpVersion      string
}

// checks for version, if new ver aval, asks user if you wanna download
func checkMCPDownload() tea.Cmd {
	return func() tea.Msg {
		cVersion, _ := getCurrentVersion()
		lRelease, errRel := getLatestDownload()
		if errRel != nil {
			return fatalError{err: errRel}
		}
		return versionCheck{CurrentVersion: cVersion, UpVersion: lRelease.Version}
	}
}

func downloadMCPCmd(progressChan chan float64) tea.Cmd {
	return func() tea.Msg {
		release, err := getLatestDownload()
		if err != nil {
			return fatalError{err: err}
		}

		downloadURL, ok := release.AssetsMap[getDownloadFile()]
		if !ok {
			return fatalError{err: fmt.Errorf("Failed to find download file for %s or '%s'", getOS(), getDownloadFile())}
		}

		client := &http.Client{Timeout: 5 * time.Minute}
		resp, err := client.Get(downloadURL)
		if err != nil {
			return fatalError{err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fatalError{err: fmt.Errorf("bad status: %s", resp.Status)}
		}

		//create file just incase
		out, err := os.Create(getFileOut())
		if err != nil {
			return fatalError{err: err}
		}
		defer out.Close()

		pw := &progressWriter{
			file:  out, // Feed it the open file descriptor
			total: resp.ContentLength,
			onProgress: func(percent float64) {
				progressChan <- percent
			},
		}

		// This cannot be optimized out; every byte MUST flow through pw.Write()
		_, err = io.Copy(pw, resp.Body)
		if err != nil {
			return fatalError{err: err}
		}

		_ = out.Sync()
		time.Sleep(250 * time.Millisecond) //small wait for progress to finish
		return downloadFinishedMsg{}
	}
}

type downloadProgressMsg float64
type downloadFinishedMsg struct{}

// progressWriter counts the bytes written to disk and sends updates
type progressWriter struct {
	file       *os.File
	total      int64
	downloaded int64
	onProgress func(float64)
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n, err := pw.file.Write(p)
	if err != nil {
		return n, err
	}

	// 2. Count the written bytes
	pw.downloaded += int64(n)
	if pw.total > 0 {
		percentage := float64(pw.downloaded) / float64(pw.total)
		pw.onProgress(percentage)
	}

	return n, nil
}
