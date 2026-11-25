package youtube

import (
	"context"
	"fmt"
	"strings"

	"github.com/ChezyName/YouTube-MCP/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func GetPublicChannel(ChannelHandle string) (ChannelStats, error) {
	if !strings.Contains(ChannelHandle, "@") {
		return ChannelStats{}, fmt.Errorf("%s needs to include an '@', Make sure this is the channel handle", ChannelHandle)
	}

	ctx := context.Background()
	svc, err := youtube.NewService(ctx, option.WithAPIKey(config.GetConfig().YouTubeAPI))
	if err != nil {
		return ChannelStats{}, err
	}

	res, err := svc.Channels.List([]string{"snippet", "statistics", "brandingSettings", "contentDetails"}).
		ForHandle(ChannelHandle).
		Do()
	if err != nil {
		return ChannelStats{}, err
	}

	if len(res.Items) == 0 {
		return ChannelStats{}, err
	}

	item := res.Items[0]

	banner := ""
	if item.BrandingSettings.Image != nil {
		banner = item.BrandingSettings.Image.BannerExternalUrl
	}

	channel := ChannelStats{
		ID:              item.Id,
		Title:           item.Snippet.Title,
		Description:     item.Snippet.Description,
		CustomURL:       item.Snippet.CustomUrl,
		PublishedAt:     item.Snippet.PublishedAt,
		Thumbnail:       item.Snippet.Thumbnails.Medium.Url,
		Banner:          banner,
		SubscriberCount: item.Statistics.SubscriberCount,
		VideoCount:      item.Statistics.VideoCount,
		TotalViewCount:  item.Statistics.ViewCount,
		Country:         item.Snippet.Country,
		Keywords:        item.BrandingSettings.Channel.Keywords,
	}

	return channel, nil
}
