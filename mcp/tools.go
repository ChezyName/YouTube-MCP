package mcp

import (
	"context"

	"github.com/ChezyName/YouTube-MCP/youtube"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "all_video",
		Title:       "Gets All Videos for User",
		Description: "Gets all public videos for user along with the data which includes: ID, Titel, Description, Thumbnail, PulishedAt",
	}, ListVideos)
}

func ListVideos(ctx context.Context, req *mcp.CallToolRequest, input interface{}) (
	*mcp.CallToolResult,
	VideoList,
	error,
) {
	videos, err := youtube.ListVideos()
	return nil, VideoList{Videos: videos, Length: len(videos)}, err
}
