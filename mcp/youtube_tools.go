package mcp

import (
	"context"
	"sync"

	"github.com/ChezyName/YouTube-MCP/youtube"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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
