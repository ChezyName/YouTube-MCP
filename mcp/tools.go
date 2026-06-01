package mcp

import (
	"context"
	"sync"

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
}

func ListVideos(ctx context.Context, req *mcp.CallToolRequest, input ListVideoParams) (
	*mcp.CallToolResult,
	VideoList,
	error,
) {
	videos, err := youtube.ListVideos(input.VideoType)
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
	analytics, err := youtube.GetChannelAnalytics(input.Range.Start(), input.Range.End())
	return nil, analytics, err
}

func GetVideoAnalytics(ctx context.Context, req *mcp.CallToolRequest, input VideoAnalyticsParams) (
	*mcp.CallToolResult,
	youtube.AnalyticsResponse,
	error,
) {
	analytics, err := youtube.GetAnalyticsForVideo(input.ID, input.Range.Start(), input.Range.End())
	return nil, analytics, err
}

// TODO: Search for video by tag, keywords, title, description, basically youtube SearchForVideo
// @returns: array of videos with basic info -> id, title, description
func SearchForVideo(ctx context.Context, req *mcp.CallToolRequest, input VideoSearchParams) (
	*mcp.CallToolResult,
	VideoSearchOutput,
	error,
) {
	var Limit = 10
	if input.Limit != nil {
		Limit = *input.Limit
	}

	basic, detailed, err := youtube.SearchVideos(input.Query, int64(Limit), input.VideoDetails, input.SearchSelf)
	return nil, VideoSearchOutput{Videos: basic, VideosDetailed: detailed}, err
}

func GetVideoTranscript(ctx context.Context, req *mcp.CallToolRequest, input VideoParams) (
	*mcp.CallToolResult,
	Transcript,
	error,
) {
	transcript, err := youtube.GetVideoTranscript(input.ID)
	if err != nil {
		return nil, Transcript{}, err
	}

	outTranscript := Transcript{}
	outTranscript.Language = transcript.Language
	outTranscript.LanguageCode = transcript.LanguageCode
	for _, snippet := range transcript.Snippets {
		outTranscript.Snippets = append(outTranscript.Snippets, TranscriptSnippet{
			Text:     snippet.Text,
			Start:    snippet.Start,
			Duration: snippet.Duration,
		})
	}

	return nil, outTranscript, err
}

func GetTopVideos(ctx context.Context, req *mcp.CallToolRequest, input TopVideosParams) (
	*mcp.CallToolResult,
	TopVideos,
	error,
) {
	var Limit = 10
	if input.Limit != nil {
		Limit = *input.Limit
	}

	videos, err := youtube.GetTopVideoIDs(input.Range.Start(), input.Range.End(), int64(Limit), input.VideoType)
	if !input.VideoDetails || err != nil {
		return nil, TopVideos{Count: len(videos), Videos: videos, Details: false}, err
	}

	var wg sync.WaitGroup
	var videosDetails = make([]*youtube.VideoDetail, len(videos))
	for i, id := range videos {
		wg.Add(1)
		go func(identifier string, index int) {
			defer wg.Done()
			details, err := youtube.GetVideo(identifier)
			if err != nil {
				videosDetails[index] = nil
			} else {
				videosDetails[index] = &details
			}

		}(id, i)
	}
	wg.Wait()

	return nil, TopVideos{Count: len(videos), Videos: videos, Details: true, VideoDetails: videosDetails}, err
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
