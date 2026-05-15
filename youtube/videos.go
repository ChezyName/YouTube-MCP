package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"context"

	"github.com/ChezyName/YouTube-MCP/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Checks if a video is a short or not
func isShort(videoID string) bool {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 2 * time.Second,
	}

	resp, err := client.Head("https://www.youtube.com/shorts/" + videoID)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// 200 means it's a Short. 301/302/303 means it's redirecting to a standard video.
	return resp.StatusCode == http.StatusOK
}

// Helper function to get Dislikes of a Video
func fetchDislikes(videoID string) (*VideoDislike, error) {
	url := fmt.Sprintf("https://returnyoutubedislikeapi.com/votes?videoId=%s", videoID)
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var d VideoDislike
	if err := json.NewDecoder(res.Body).Decode(&d); err != nil {
		return nil, err
	}
	return &d, nil
}

func ListVideos(vType *VideoType) ([]Video, error) {
	ctx := context.Background()

	if config.GetConfig().ChannelHandle == "" {
		//log.Error(`Empty Channel Handle, Update .env with ChannelHandle="YOUR_CHANNEL_HANDLE"`)
		return []Video{}, fmt.Errorf("Missing Channel Handle")
	}

	svc, err := youtube.NewService(ctx, option.WithAPIKey(config.GetConfig().YouTubeAPI))
	if err != nil {
		//log.Error("Failed to create YouTube client")
		return []Video{}, err
	}

	//Gets User (owner of the API key)'s Channel
	channelRes, err := svc.Channels.List([]string{"contentDetails"}).ForHandle(config.GetConfig().ChannelHandle).Do()
	if err != nil {
		//log.Error("Failed to fetch channel: " + err.Error())
		return []Video{}, err
	}

	uploadsPlaylistID := channelRes.Items[0].ContentDetails.RelatedPlaylists.Uploads

	var videos []Video
	nextPageToken := ""

	type result struct {
		index int
		vType VideoType
	}
	resultsChan := make(chan result)
	var wg sync.WaitGroup

	for {
		call := svc.PlaylistItems.List([]string{"snippet"}).
			PlaylistId(uploadsPlaylistID).
			MaxResults(50).
			PageToken(nextPageToken)

		res, err := call.Do()
		if err != nil {
			//log.Error("Failed to fetch videos: " + err.Error())
			return []Video{}, err
		}

		//TODO: add mutex to async get the short type
		for _, item := range res.Items {
			videos = append(videos, Video{
				ID:          item.Snippet.ResourceId.VideoId,
				Title:       item.Snippet.Title,
				Description: item.Snippet.Description,
				PublishedAt: item.Snippet.PublishedAt,
				Thumbnail:   item.Snippet.Thumbnails.Medium.Url,
				Type:        Unknown,
			})
			var idx = len(videos) - 1

			wg.Add(1)
			go func(i int, vidID string) {
				defer wg.Done()

				videoType := Longform
				if isShort(vidID) {
					videoType = Short
				}

				resultsChan <- result{index: i, vType: videoType}
			}(idx, item.Snippet.ResourceId.VideoId)
		}

		nextPageToken = res.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	for res := range resultsChan {
		videos[res.index].Type = res.vType
	}

	//Filter quickly
	n := 0
	for _, v := range videos {
		// If user didn't specify (Both/nil) OR if it matches exactly
		if vType == nil || *vType == Both || v.Type == *vType {
			videos[n] = v
			n++
		}
	}
	videos = videos[:n]
	return videos, nil
}

func GetVideo(videoID string) (VideoDetail, error) {
	ctx := context.Background()

	if videoID == "" {
		//log.Error("Missing video ID", http.StatusBadRequest)
		return VideoDetail{}, fmt.Errorf("Missing Video ID")
	}

	svc, err := youtube.NewService(ctx, option.WithAPIKey(config.GetConfig().YouTubeAPI))
	if err != nil {
		//log.Error("Failed to create YouTube client")
		return VideoDetail{}, err
	}

	res, err := svc.Videos.List([]string{"snippet", "statistics", "contentDetails"}).Id(videoID).Do()
	if err != nil {
		//log.Error("Failed to fetch video: " + err.Error())
		return VideoDetail{}, err
	}

	if len(res.Items) == 0 {
		//log.Error("Video not found", http.StatusNotFound)
		return VideoDetail{}, err
	}

	item := res.Items[0]
	//dislikes removed, use the dislikes API [returnyoutubedislikeapi]
	dislikes, err := fetchDislikes(item.Id)
	var total_dislikes uint64 = 0

	if err != nil {
		fmt.Println("Error getting dislikes from Return YouTube Dislike API")
	} else {
		total_dislikes = uint64(dislikes.Dislikes)
	}

	videoType := Longform
	if isShort(videoID) {
		videoType = Short
	}

	video := VideoDetail{
		ID:           item.Id,
		Title:        item.Snippet.Title,
		Description:  item.Snippet.Description,
		PublishedAt:  item.Snippet.PublishedAt,
		Thumbnail:    item.Snippet.Thumbnails.Medium.Url,
		Duration:     item.ContentDetails.Duration, // ISO 8601 e.g. PT4M
		ViewCount:    item.Statistics.ViewCount,
		DislikeCount: total_dislikes,
		LikeCount:    item.Statistics.LikeCount,
		CommentCount: item.Statistics.CommentCount,
		Type:         videoType,
	}

	return video, nil
}

func GetVideoComments(videoID string, limit int) (CommentsResponse, error) {
	ctx := context.Background()

	svc, err := youtube.NewService(ctx, option.WithAPIKey(config.GetConfig().YouTubeAPI))
	if err != nil {
		//log.Error("Failed to create YouTube client")
		return CommentsResponse{}, err
	}

	res, err := svc.CommentThreads.List([]string{"snippet"}).
		VideoId(videoID).
		Order("relevance"). // top comments
		MaxResults(int64(limit)).
		Do()
	if err != nil {
		//log.Error("Failed to fetch comments: " + err.Error())
		return CommentsResponse{}, err
	}

	var comments []Comment
	for _, item := range res.Items {
		c := item.Snippet.TopLevelComment.Snippet
		comments = append(comments, Comment{
			ID:          item.Id,
			Author:      c.AuthorDisplayName,
			Text:        c.TextDisplay,
			LikeCount:   c.LikeCount,
			PublishedAt: c.PublishedAt,
			UpdatedAt:   c.UpdatedAt,
		})
	}

	outComment := CommentsResponse{
		VideoID:  videoID,
		Total:    len(comments),
		Limit:    limit,
		Comments: comments,
	}

	return outComment, nil
}
