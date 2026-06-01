package youtube

import (
	"context"
	"net/http"
	"time"

	"github.com/ChezyName/YouTube-MCP/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func AuthCheck() bool {
	cfg := config.GetConfig()
	if cfg.YouTubeRefreshToken == "" || (cfg.YOUTUBE_CLIENT_ID == "" || cfg.YOUTUBE_CLIENT_SECRET == "") {
		return false
	}

	//test to see if end-point can be reached
	client, err := config.GetOAuthClient()
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// We use an authenticated request to the YouTube metadata validation API.
	// Filtering by mine=true uses almost zero quota but requires operational OAuth credentials.
	req, err := http.NewRequestWithContext(ctx, "GET",
		"https://www.googleapis.com/youtube/v3/channels?part=snippet,brandingSettings&mine=true", nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func APICheck() bool {
	cfg := config.GetConfig()
	if cfg == nil {
		return false
	}

	if cfg.YouTubeAPI == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	svc, err := youtube.NewService(ctx, option.WithAPIKey(cfg.YouTubeAPI))
	if err != nil {
		return false
	}

	_, err = svc.VideoCategories.List([]string{"snippet"}).RegionCode("US").Do()

	return err == nil
}
