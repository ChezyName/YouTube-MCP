package mcp

import (
	"github.com/ChezyName/YouTube-MCP/youtube"
)

type VideoList struct {
	Length int             `json:"length" jsonschema:"The number of videos"`
	Videos []youtube.Video `json:"videos" jsonschema:"The videos"`
}

type VideoParams struct {
	ID string `json:"id" jsonschema:"YouTube video ID"`
}

type TopVideosParams struct {
	Range        *youtube.Range     `json:"range,omitempty" jsonschema:"Optional: Date range in 'YYYY-MM-DD/YYYY-MM-DD' format, 'lifetime', or number of days (e.g. '28'). Defaults to 30 days."`
	Limit        *int               `json:"limit,omitempty" jsonschema:"The max number of videos to return"`
	VideoDetails bool               `json:"details" jsonschema:"Boolean which dictates grabbing the video details"`
	VideoType    *youtube.VideoType `json:"content_type,omitempty" jsonschema:"The type of content, 'short', 'long', or 'both' to filter through, defaults to both"`
}

type VideoCommentsParams struct {
	ID    string `json:"id" jsonschema:"YouTube video ID"`
	Limit *int   `json:"limit,omitempty" jsonschema:"The max number of comments loaded. default=20"`
}

type ChannelAnalyticsParams struct {
	Range *youtube.Range `json:"range,omitempty" jsonschema:"Optional: Date range in 'YYYY-MM-DD/YYYY-MM-DD' format, 'lifetime', or number of days (e.g. '28'). Defaults to 30 days."`
}

type VideoAnalyticsParams struct {
	ID    string         `json:"id" jsonschema:"YouTube video ID"`
	Range *youtube.Range `json:"range,omitempty" jsonschema:"Optional: Date range in 'YYYY-MM-DD/YYYY-MM-DD' format, 'lifetime', or number of days (e.g. '28'). Defaults to 30 days."`
}

type ListVideoParams struct {
	VideoType *youtube.VideoType `json:"content_type,omitempty" jsonschema:"The type of content, 'short', 'long', or 'both' to filter through, defaults to both"`
}

type VideoSearchParams struct {
	Query        string `json:"query" jsonschema:"The query string that it uses to search for videos"`
	Limit        *int   `json:"limit,omitempty" jsonschema:"The max number of comments loaded. default=20"`
	VideoDetails bool   `json:"details" jsonschema:"Boolean which dictates grabbing the video details"`
	SearchSelf   bool   `json:"search_self" jsonschema:"Boolean which dictates if search is self-channel, or all of youtube"`
}

type VideoSearchOutput struct {
	Videos         []*youtube.Video       `json:"videos" jsonschema:"The found videos"`
	VideosDetailed []*youtube.VideoDetail `json:"videos_detailed" jsonschema:"The found videos with details"`
}

type TranscriptSnippet struct {
	Text     string  `json:"text"`
	Start    float64 `json:"start"`
	Duration float64 `json:"duration"`
}

type Transcript struct {
	Snippets     []TranscriptSnippet `json:"snippets"`
	Language     string              `json:"language"`
	LanguageCode string              `json:"language_code"`
}

type TopVideos struct {
	Count        int                    `json:"count" jsonschema:"The number of videos found"`
	Videos       []string               `json:"videos" jsonschema:"The video IDs"`
	Details      bool                   `json:"details" jsonschema:"If video details are avalable"`
	VideoDetails []*youtube.VideoDetail `json:"video_details" jsonschema:"The detailed list of top x videos"`
}

type AuthCheckResult struct {
	IsAuthenticatedAnalytics bool `json:"authenticated_analytics" jsonschema:"If the Private Analytics Data is Authenticated"`
	IsAuthenticatedData      bool `json:"authenticated_data" jsonschema:"If the Public Data API v3 is Authenticated"`
}
