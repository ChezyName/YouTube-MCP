package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ChezyName/YouTube-MCP/config"
	youtubemcp "github.com/ChezyName/YouTube-MCP/mcp"
	"github.com/ChezyName/YouTube-MCP/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// main.go
var Version = "UNKNOWN" // default if not built with ldflags

func init() {
	config.LoadConfig()

	//init competitors if allowed
	if tools.IsCompetitorsEnabled() {
		tools.GetCompetitors() //loads the file
	}
}

func main() {
	//return version upon `./YouTube-MCP -v`
	if len(os.Args) > 1 && os.Args[1] == "-v" {
		fmt.Println(Version)
		os.Exit(0)
	}

	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{
		Name:       "YouTube MCP",
		Title:      "YouTube MCP",
		Version:    Version,
		WebsiteURL: "https://github.com/ChezyName/YouTube-MCP",
	}, nil)
	youtubemcp.AddTools(server)

	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
