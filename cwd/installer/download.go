package main

import (
	"path/filepath"
	"runtime"

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
		// Windows only needs amd64 based on your build script
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

func getCurrentVersion() {

}

type downloadFinishedMsg struct{ path string }

// checks for version, if new ver aval, asks user if you wanna download
func checkMCPDownload() {

}

func downloadMCPCmd() tea.Cmd {
	return func() tea.Msg {
		// download logic...
		return downloadFinishedMsg{path: "/path/to/mcp"}
	}
}
