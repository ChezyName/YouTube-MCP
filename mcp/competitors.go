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

func AddCompetitor(ctx context.Context, req *mcp.CallToolRequest, input CompitorInput) (
	*mcp.CallToolResult,
	interface{},
	error,
) {
	var tags = []string{}
	if input.Tags != nil {
		tags = *input.Tags
	}
	err := tools.AddCompetitor(input.Name, tags)
	return nil, nil, err
}

func RemoveCompetitor(ctx context.Context, req *mcp.CallToolRequest, input CompitorRemoveInput) (
	*mcp.CallToolResult,
	interface{},
	error,
) {
	err := tools.RemoveCompetitor(input.Name)
	return nil, nil, err
}
