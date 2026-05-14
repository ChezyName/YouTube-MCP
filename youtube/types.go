package youtube

type VideoDislike struct {
	Likes       int     `json:"likes"`
	Dislikes    int     `json:"dislikes"`
	RawLikes    int     `json:"rawLikes"`
	RawDislikes int     `json:"rawDislikes"`
	Rating      float64 `json:"rating"`
	ViewCount   int     `json:"viewCount"`
	Deleted     bool    `json:"deleted"`
}

type Video struct {
	ID          string `json:"id" jsonschema:"The unique ID of the video"`
	Title       string `json:"title" jsonschema:"The title of the video"`
	Description string `json:"description" jsonschema:"The description of the video"`
	PublishedAt string `json:"published_at" jsonschema:"The time in which the vieo was published"`
	Thumbnail   string `json:"thumbnail" jsonschema:"The thumbnail url of the video"`
}

type VideoDetail struct {
	ID           string `json:"id" jsonschema:"The unique ID of the video"`
	Title        string `json:"title" jsonschema:"The title of the video"`
	Description  string `json:"description" jsonschema:"The description of the video"`
	PublishedAt  string `json:"published_at" jsonschema:"The time in which the vieo was published"`
	Thumbnail    string `json:"thumbnail" jsonschema:"The thumbnail url of the video"`
	Duration     string `json:"duration" jsonschema:"The duration of the video"`
	ViewCount    uint64 `json:"view_count" jsonschema:"The number of video views over the video's lifetime"`
	DislikeCount uint64 `json:"disview_count" jsonschema:"The number of dislikes over the video's lifetime (using return YouTube Dislikes)"`
	LikeCount    uint64 `json:"like_count" jsonschema:"The number of likes over the video's lifetime"`
	CommentCount uint64 `json:"comment_count" jsonschema:"The number of comments over the video's lifetime"`
}
type DateRange struct {
	Start string `json:"start" jsonschema:"start date"`
	End   string `json:"end" jsonschema:"end date"`
}

type OverviewStats struct {
	Views          float64 `json:"views" jsonschema:"Video views"`
	WatchTimeHours float64 `json:"watch_time_hours" jsonschema:"Video watch time in hours"`
	AVD            float64 `json:"avg_view_duration_seconds" jsonschema:"Video average view duration"`
	AVP            float64 `json:"avg_view_percentage" jsonschema:"Video average view percentage"`
}

type EngagementStats struct {
	Likes       float64 `json:"likes" jsonschema:"Video likes"`
	Dislikes    float64 `json:"dislikes" jsonschema:"Video dislikes (return YouTube Dislikes)"`
	Comments    float64 `json:"comments" jsonschema:"Number of Video Comments"`
	Shares      float64 `json:"shares" jsonschema:"Number of Video Shares"`
	Subscribers float64 `json:"subscribers_gained" jsonschema:"Number of Subscribers Gained on this Video"`
}

type ImpressionStats struct {
	Impressions float64 `json:"impressions" jsonschema:"Number of Impressions"`
	CTR         float64 `json:"click_through_rate" jsonschema:"Click through Rate"`
}

type RowData struct {
	Label  string  `json:"label" jsonschema:"title"`
	Value  float64 `json:"value" jsonschema:"value"`
	Value2 float64 `json:"value2,omitempty"` // for multi-metric rows e.g. retention %
}

type AnalyticsResponse struct {
	VideoID     string          `json:"video_id" jsonschema:"The ID of the video"`
	DateRange   DateRange       `json:"date_range" jsonschema:"Base stats such as views, watch time, AVD and AVP"`
	Overview    OverviewStats   `json:"overview" jsonschema:"Base stats such as views, watch time, AVD and AVP"`
	Engagement  EngagementStats `json:"engagement" jsonschema:"Engagement stats for the video"`
	Impressions ImpressionStats `json:"impressions" jsonschema:"Impression stats for the video"`
	//Audience       AudienceStats   `json:"audience"`
	TrafficSources []RowData `json:"traffic_sources" jsonschema:"Where users are seeing the video"`
	Retention      []RowData `json:"retention" jsonschema:"The retention graph"`
	Geography      []RowData `json:"geography" jsonschema:"Where users are watching the videos"`
	DeviceTypes    []RowData `json:"device_types" jsonschema:"The devices users are watching the video on"`
	DailyBreakdown []RowData `json:"daily_breakdown" jsonschema:"Daily stats"`
}

type ChannelStats struct {
	ID              string `json:"id" jsonschema:"Channel ID"`
	Title           string `json:"title" jsonschema:"Title or name of the channel"`
	Description     string `json:"description" jsonschema:"Channel Description"`
	CustomURL       string `json:"custom_url" jsonschema:"Channel URL - typically the Handle"`
	PublishedAt     string `json:"published_at" jsonschema:"When the channel was publisehd or created"`
	Thumbnail       string `json:"thumbnail" jsonschema:"The channel icon url"`
	Banner          string `json:"banner" jsonschema:"The channel banner url"`
	SubscriberCount uint64 `json:"subscriber_count" jsonschema:"The number of subscribers"`
	VideoCount      uint64 `json:"video_count" jsonschema:"The number of videos"`
	TotalViewCount  uint64 `json:"total_view_count" jsonschema:"The total views"`
	Country         string `json:"country" jsonschema:"The country"`
}

type ChannelAnalyticsResponse struct {
	DateRange        DateRange     `json:"date_range" jsonschema:"Date range for the channel data"`
	Overview         OverviewStats `json:"overview" jsonschema:"Base stats such as views, watch time, AVD and AVP"`
	SubscriberGrowth []RowData     `json:"subscriber_growth" jsonschema:"Graph of Subscriber growth"`
	TopVideos        []VideoDetail `json:"top_videos" jsonschema:"Top 10 views"`
	TrafficSources   []RowData     `json:"traffic_sources" jsonschema:"Where users are seeing the video from"`
	Geography        []RowData     `json:"geography" jsonschema:"Locations of those who watch the video"`
	DeviceTypes      []RowData     `json:"device_types" jsonschema:"What devices users use"`
	AgeGroups        []RowData     `json:"age_groups" jsonschema:"The age range of users"`
	Gender           []RowData     `json:"gender" jsonschema:"The gender of users"`
	DailyBreakdown   []RowData     `json:"daily_breakdown" jsonschema:"Breakdown of views, subs gained, likes, shares per day"`
}

type Comment struct {
	ID          string `json:"id" jsonschema:"The ID of the comment"`
	Author      string `json:"author" jsonschema:"The author of the comment"`
	Text        string `json:"text" jsonschema:"The text of the comment"`
	LikeCount   int64  `json:"like_count" jsonschema:"The number of likes of the comment"`
	PublishedAt string `json:"published_at" jsonschema:"When the comment was published"`
	UpdatedAt   string `json:"updated_at" jsonschema:"When the comment was updated"`
}

type CommentsResponse struct {
	VideoID  string    `json:"video_id" jsonschema:"The associated video"`
	Total    int       `json:"total" jsonschema:"The total number of comments retrieved"`
	Limit    int       `json:"limit" jsonschema:"The given limit"`
	Comments []Comment `json:"comments" jsonschema:"The comments in the video"`
}
