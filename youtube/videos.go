package youtube

import (
	"encoding/json"
	"fmt"
	"net/http"

	"context"

	"github.com/ChezyName/YouTube-MCP/config"
	"github.com/gorilla/mux"
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

func ListVideos(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	svc, err := youtube.NewService(ctx, option.WithAPIKey(config.GetConfig().YouTubeAPI))
	if err != nil {
		http.Error(w, "Failed to create YouTube client", http.StatusInternalServerError)
		return
	}

	//Gets User (owner of the API key)'s Channel
	channelRes, err := svc.Channels.List([]string{"contentDetails"}).ForHandle(config.GetConfig().ChannelHandle).Do()
	if err != nil {
		http.Error(w, "Failed to fetch channel: "+err.Error(), http.StatusInternalServerError)
		return
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
			http.Error(w, "Failed to fetch videos: "+err.Error(), http.StatusInternalServerError)
			return
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videos)
}

func GetVideo(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	vars := mux.Vars(r)
	videoID := vars["id"]

	if videoID == "" {
		fmt.Println("Missing Video ID for GetVideo")
		http.Error(w, "Missing video ID", http.StatusBadRequest)
		return
	}

	svc, err := youtube.NewService(ctx, option.WithAPIKey(config.GetConfig().YouTubeAPI))
	if err != nil {
		fmt.Println("Failed to create YouTube API Client")
		http.Error(w, "Failed to create YouTube client", http.StatusInternalServerError)
		return
	}

	res, err := svc.Videos.List([]string{"snippet", "statistics", "contentDetails"}).Id(videoID).Do()
	if err != nil {
		fmt.Println("Failed to fetch video: " + err.Error())
		http.Error(w, "Failed to fetch video: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(res.Items) == 0 {
		fmt.Println("Video was not found")
		http.Error(w, "Video not found", http.StatusNotFound)
		return
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(video)
}
