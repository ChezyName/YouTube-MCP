package mcp

type ToolParam struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type InputSchema struct {
	Type       string               `json:"type"`
	Properties map[string]ToolParam `json:"properties"`
	Required   []string             `json:"required,omitempty"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema InputSchema `json:"inputSchema"`
}

func GetTools() []Tool {
	return []Tool{
		{
			Name:        "list_videos",
			Description: "List all public and unlisted videos on the YouTube channel",
			InputSchema: InputSchema{
				Type:       "object",
				Properties: map[string]ToolParam{},
			},
		},
		{
			Name:        "get_video",
			Description: "Get details for a specific YouTube video including stats, duration, likes, dislikes, and rating",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]ToolParam{
					"id": {Type: "string", Description: "YouTube video ID"},
				},
				Required: []string{"id"},
			},
		},
		{
			Name:        "get_video_analytics",
			Description: "Get full analytics for a specific video including views, watch time, CTR, AVD, retention curve, traffic sources, geography, and device types",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]ToolParam{
					"id":    {Type: "string", Description: "YouTube video ID"},
					"range": {Type: "string", Description: "Date range: number of days (e.g. 90, 365) or 'lifetime'"},
				},
				Required: []string{"id"},
			},
		},
		{
			Name:        "get_video_comments",
			Description: "Get top comments for a specific YouTube video",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]ToolParam{
					"id":    {Type: "string", Description: "YouTube video ID"},
					"limit": {Type: "string", Description: "Number of comments to return (default 20, max 100)"},
				},
				Required: []string{"id"},
			},
		},
		{
			Name:        "get_channel",
			Description: "Get public channel stats including subscriber count, total views, video count, and branding info",
			InputSchema: InputSchema{
				Type:       "object",
				Properties: map[string]ToolParam{},
			},
		},
		{
			Name:        "get_channel_analytics",
			Description: "Get channel-wide analytics including views, watch time, CTR, subscriber growth, top videos, traffic sources, geography, age groups, and gender breakdown",
			InputSchema: InputSchema{
				Type: "object",
				Properties: map[string]ToolParam{
					"range": {Type: "string", Description: "Date range: number of days (e.g. 90, 365) or 'lifetime'"},
				},
			},
		},
	}
}
