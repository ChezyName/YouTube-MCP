package mcp

import (
	"context"

	"github.com/ChezyName/YouTube-MCP/tools"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func GetCompetitorsList(ctx context.Context, req *mcp.CallToolRequest, input interface{}) (
	*mcp.CallToolResult,
	GetCompetitorsListOutput,
	error,
) {
	data, err := tools.GetCompetitors()
	return nil, GetCompetitorsListOutput{Competitors: data}, err
}
