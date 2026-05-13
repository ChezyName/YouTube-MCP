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
