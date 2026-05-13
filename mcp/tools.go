package mcp

import (
	"context"

	"github.com/ChezyName/YouTube-MCP/youtube"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func AddTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "all_videos",
		Title:       "Gets All Videos",
		Description: "[PUBLIC API] Gets all public videos for user along with the data which includes: ID, Title, Description, Thumbnail, PulishedAt",
	}, ListVideos)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_video",
		Title:       "Gets Single Video Detailed",
		Description: "[PUBLIC API] Gets a detailed of a single video with such as: ID, Title, Description, Thumbnail, PulishedAt, Duration, Views, Dislikes, Likes, CommentCount",
	}, GetVideo)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_video_comments",
		Title:       "Gets Comments from a Video",
		Description: "[PUBLIC API] Gets a limited number of coments from a given video, returns the author, text, likecount, publishedat, updatedat, and id of the comment",
	}, GetVideoComments)
}

func ListVideos(ctx context.Context, req *mcp.CallToolRequest, input interface{}) (
	*mcp.CallToolResult,
	VideoList,
	error,
) {
	videos, err := youtube.ListVideos()
	return nil, VideoList{Videos: videos, Length: len(videos)}, err
}

func GetVideo(ctx context.Context, req *mcp.CallToolRequest, input VideoParams) (
	*mcp.CallToolResult,
	youtube.VideoDetail,
	error,
) {
	video, err := youtube.GetVideo(input.ID)
	return nil, video, err
}

func GetVideoComments(ctx context.Context, req *mcp.CallToolRequest, input VideoCommentsParams) (
	*mcp.CallToolResult,
	youtube.CommentsResponse,
	error,
) {
	limit := 20 // Set our default
	if input.Limit != nil {
		limit = *input.Limit
	}

	comments, err := youtube.GetVideoComments(input.ID, limit)
	return nil, comments, err
}
