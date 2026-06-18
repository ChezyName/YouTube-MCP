package mcp

import (
	"context"

	"github.com/ChezyName/YouTube-MCP/tools"
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
		Description: "[PUBLIC API] Gets the users channel with only the public information such as name, handle, bannel, icon, and description.",
	}, GetChannel)

	//Analytics from the private youtube-analytics API
	mcp.AddTool(server, &mcp.Tool{
		Name:  "get_channel_analytics",
		Title: "Channel Analytics Data",
		Description: `[PRIVATE API] Get the users channel Analytics given a start_date and end_date (or range in days), defaults to lifetime
		Returns views, watch time hours, average view duration, average view percentage, likes, dislikes, comments, shares, subscribers, impressions,
		click through rate, unique views, subscriber grouth graph, top views, traffic sources, geography, device types, age groups, gender, and daily breakdown`,
	}, GetChannelAnalytics)

	mcp.AddTool(server, &mcp.Tool{
		Name:  "get_video_analytics",
		Title: "Video Analytics Data",
		Description: `[PRIVATE API] Get the users Video Analytics given a video_id, start_date and end_date (or range in days), defaults to last 90 days,
		Returns views, watch time hours, average view duration, average view percentage, likes, dislikes, comments, shares, subscribers, impressions,
		click through rate, unique views, traffic sources, geography, device types, age groups, gender, and daily breakdown`,
	}, GetVideoAnalytics)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_video_transcript",
		Title:       "Video Transcript",
		Description: "[PUBLIC API] Returns a structured list of the video transcript",
	}, GetVideoTranscript)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_top_videos",
		Title:       "Top Videos",
		Description: "[PUBLIC API] Returns a list of the top videos given a date range. Defaults to last 90days and top 10. Details show the public basic data such as name, description, and likes.",
	}, GetTopVideos)

	mcp.AddTool(server, &mcp.Tool{
		Name:  "search_video",
		Title: "Search for Videos given a query",
		Description: `[PUBLIC API] Returns a detailed or basic list of videos found when searching for them.
		Basic: ID, Title, Description, Thumbnail, PulishedAt
		Detailed: ID, Title, Description, Thumbnail, PulishedAt, Duration, Views, Dislikes, Likes, CommentCount
		`,
	}, SearchForVideo)

	//Auth Check / Test Tools
	mcp.AddTool(server, &mcp.Tool{
		Name:        "auth_check",
		Title:       "Authentication Check",
		Description: "Checks both YouTube Data API and Analytics API",
	}, ProgramAuthCheck)

	//Competitors - only if config allows it
	if tools.IsCompetitorsEnabled() {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "competitors_list",
			Title:       "Get Competitors",
			Description: "Returns a list of Competitors written by the User, or an AI",
		}, GetCompetitorsList)

		mcp.AddTool(server, &mcp.Tool{
			Name:        "add_competitor",
			Title:       "Add Competitor",
			Description: "Adds a Competitor to the list of competitors for the user",
		}, AddCompetitor)

		mcp.AddTool(server, &mcp.Tool{
			Name:        "remove_competitor",
			Title:       "Remove Competitor",
			Description: "Removes a Competitor to the list of competitors for the user",
		}, RemoveCompetitor)
	}
}

// Add all Auth Checks
func ProgramAuthCheck(ctx context.Context, req *mcp.CallToolRequest, input interface{}) (
	*mcp.CallToolResult,
	AuthCheckResult,
	error,
) {
	return nil, AuthCheckResult{
		IsAuthenticatedAnalytics: youtube.AuthCheck(),
		IsAuthenticatedData:      youtube.APICheck(),
	}, nil
}
