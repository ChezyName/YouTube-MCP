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
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishedAt string `json:"published_at"`
	Thumbnail   string `json:"thumbnail"`
}

type VideoDetail struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	PublishedAt  string `json:"published_at"`
	Thumbnail    string `json:"thumbnail"`
	Duration     string `json:"duration"`
	ViewCount    uint64 `json:"view_count"`
	DislikeCount uint64 `json:"disview_count"`
	LikeCount    uint64 `json:"like_count"`
	CommentCount uint64 `json:"comment_count"`
}
type DateRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type OverviewStats struct {
	Views          float64 `json:"views"`
	WatchTimeHours float64 `json:"watch_time_hours"`
	AVD            float64 `json:"avg_view_duration_seconds"`
	AVP            float64 `json:"avg_view_percentage"`
}

type EngagementStats struct {
	Likes       float64 `json:"likes"`
	Dislikes    float64 `json:"dislikes"`
	Comments    float64 `json:"comments"`
	Shares      float64 `json:"shares"`
	Subscribers float64 `json:"subscribers_gained"`
}

type ImpressionStats struct {
	Impressions float64 `json:"impressions"`
	CTR         float64 `json:"click_through_rate"`
	UniqueViews float64 `json:"unique_viewers"`
}

type RowData struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Value2 float64 `json:"value2,omitempty"` // for multi-metric rows e.g. retention %
}

type AnalyticsResponse struct {
	VideoID     string          `json:"video_id"`
	DateRange   DateRange       `json:"date_range"`
	Overview    OverviewStats   `json:"overview"`
	Engagement  EngagementStats `json:"engagement"`
	Impressions ImpressionStats `json:"impressions"`
	//Audience       AudienceStats   `json:"audience"`
	TrafficSources []RowData `json:"traffic_sources"`
	Retention      []RowData `json:"retention"`
	Geography      []RowData `json:"geography"`
	DeviceTypes    []RowData `json:"device_types"`
	DailyBreakdown []RowData `json:"daily_breakdown"`
}
