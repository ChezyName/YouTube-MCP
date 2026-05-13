package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"

	"context"

	"github.com/ChezyName/YouTube-MCP/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

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

func ListVideos() ([]Video, error) {
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

		for _, item := range res.Items {
			videos = append(videos, Video{
				ID:          item.Snippet.ResourceId.VideoId,
				Title:       item.Snippet.Title,
				Description: item.Snippet.Description,
				PublishedAt: item.Snippet.PublishedAt,
				Thumbnail:   item.Snippet.Thumbnails.Medium.Url,
			})
		}

		nextPageToken = res.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

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
