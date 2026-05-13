package youtube

import (
	"encoding/json"
	"net/http"

	"context"

	"github.com/ChezyName/YouTube-MCP/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

type Video struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PublishedAt string `json:"published_at"`
	Thumbnail   string `json:"thumbnail"`
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

	// Step 2: Fetch videos from uploads playlist
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
