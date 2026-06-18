package tools

type Competitor struct {
	Name string   `json:"name" jsonschema:"The name / handle of the channel"`
	Tags []string `json:"tags,omitempty" jasonschema:"Optional field for tags related to this channel, explains what channel does"`
}
