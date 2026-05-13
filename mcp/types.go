package mcp

import "github.com/ChezyName/YouTube-MCP/youtube"

type VideoList struct {
	Length int             `json:"length" jsonschema:"The number of videos"`
	Videos []youtube.Video `json:"videos" jsonschema:"The videos"`
}

type GetVideoParams struct {
	ID string `json:"id" jsonschema:"description=YouTube video ID,required"`
}
