package main

import (
	"context"
	"log"

	"github.com/ChezyName/YouTube-MCP/config"
	youtubemcp "github.com/ChezyName/YouTube-MCP/mcp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func init() {
	config.LoadConfig()
	//config.EnsureRefreshToken() - only allowed for clients
}

func main() {
	// Create a server with a single tool.
	server := mcp.NewServer(&mcp.Implementation{Name: "YouTube MCP", Version: "v1.0.0"}, nil)
	youtubemcp.AddTools(server)

	// Run the server over stdin/stdout, until the client disconnects.
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
}
