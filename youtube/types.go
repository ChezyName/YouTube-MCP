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
