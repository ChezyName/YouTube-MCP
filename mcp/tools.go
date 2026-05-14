package mcp

import (
	"context"

	"github.com/ChezyName/YouTube-MCP/youtube"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func AddTools(server *mcp.Server) {
	//Video Stats
	mcp.AddTool(server, &mcp.Tool{
		Name:        "all_videos",
		Title:       "All Public Videos",
		Description: "[PUBLIC API] Gets all public videos for user along with the data which includes: ID, Title, Description, Thumbnail, PulishedAt",
	}, ListVideos)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_video",
		Title:       "Video Details",
		Description: "[PUBLIC API] Gets a detailed of a single video with such as: ID, Title, Description, Thumbnail, PulishedAt, Duration, Views, Dislikes, Likes, CommentCount",
	}, GetVideo)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_video_comments",
		Title:       "Video Comments",
		Description: "[PUBLIC API] Gets a limited number of coments from a given video, returns the author, text, likecount, publishedat, updatedat, and id of the comment",
	}, GetVideoComments)

	//Channel Stats
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_channel",
		Title:       "Channel Data",
		Description: "[PUBLIC API] Gets a limited number of coments from a given video, returns the author, text, likecount, publishedat, updatedat, and id of the comment",
	}, GetChannel)

	mcp.AddTool(server, &mcp.Tool{
		Name:  "get_channel_analytics",
		Title: "Channel Analytics Data",
		Description: `[PRIVATE API] Get the users channel Analytics given a start_date and end_date (or range in days),
		Returns views, watch time hours, average view duration, average view percentage, likes, dislikes, comments, shares, subscribers, impressions,
		click through rate, unique views, subscriber grouth graph, top views, traffic sources, geography, device types, age groups, gender, and daily breakdown`,
	}, GetChannelAnalytics)
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

func GetChannel(ctx context.Context, req *mcp.CallToolRequest, input interface{}) (
	*mcp.CallToolResult,
	youtube.ChannelStats,
	error,
) {
	channel, err := youtube.GetChannel()
	return nil, channel, err
}

func GetChannelAnalytics(ctx context.Context, req *mcp.CallToolRequest, input ChannelAnalyticsParams) (
	*mcp.CallToolResult,
	youtube.ChannelAnalyticsResponse,
	error,
) {
	analytics, err := youtube.GetChannelAnalytics(input.StartDate, input.EndDate, input.Range)
	return nil, analytics, err
}
