package mcp

import "github.com/ChezyName/YouTube-MCP/youtube"

type VideoList struct {
	Length int             `json:"length" jsonschema:"The number of videos"`
	Videos []youtube.Video `json:"videos" jsonschema:"The videos"`
}

type VideoParams struct {
	ID string `json:"id" jsonschema:"YouTube video ID"`
}

type VideoCommentsParams struct {
	ID    string `json:"id" jsonschema:"YouTube video ID"`
	Limit *int   `json:"limit,omitempty" jsonschema:"The max number of comments loaded. default=20"`
}

type ChannelAnalyticsParams struct {
	Range     string `json:"range" jsonschema:"Range in days, overrides the start_date and end_date, lifetime is accepted"`
	StartDate string `json:"start_date" jsonschema:"The start date in YYYY-MM-DD format"`
	EndDate   string `json:"end_date" jsonschema:"The end date in YYYY-MM-DD format"`
}
