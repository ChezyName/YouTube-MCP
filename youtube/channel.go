package youtube

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sync"

	"github.com/ChezyName/YouTube-MCP/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	youtubeAnalytics "google.golang.org/api/youtubeanalytics/v2"
)

// Get Links
func extractLinks(text string) []string {
	re := regexp.MustCompile(`https?://[^\s]+`)
	return re.FindAllString(text, -1)
}

func GetChannel() (ChannelStats, error) {
	ctx := context.Background()
	svc, err := youtube.NewService(ctx, option.WithAPIKey(config.GetConfig().YouTubeAPI))
	if err != nil {
		return ChannelStats{}, err
	}

	res, err := svc.Channels.List([]string{"snippet", "statistics", "brandingSettings", "contentDetails"}).
		ForHandle(config.GetConfig().ChannelHandle).
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

/*
(defaults to lifetime)
Params:

	start: start date for range, needs end date (defaults to -90days)
	end: end date for range
	range: custom numbers, in days or Lifetime - this superseeds all
*/
func GetChannelAnalytics(startDate string, endDate string) (ChannelAnalyticsResponse, error) {
	client, err := config.GetOAuthClient()
	if err != nil {
		return ChannelAnalyticsResponse{}, err
	}

	svc, err := youtubeAnalytics.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return ChannelAnalyticsResponse{}, err
	}

	analytics := ChannelAnalyticsResponse{
		DateRange: DateRange{Start: startDate, End: endDate},
	}

	fetchChannelMetrics := func(metrics, dimensions string) (*youtubeAnalytics.QueryResponse, error) {
		call := svc.Reports.Query().
			Ids("channel==MINE").
			StartDate(startDate).
			EndDate(endDate).
			Metrics(metrics)
		if dimensions != "" {
			call = call.Dimensions(dimensions)
		}
		res, err := call.Do()
		if err != nil {
			// THIS WILL TELL YOU EXACTLY WHICH CALL FAILED
			fmt.Fprintf(os.Stderr, "[API ERROR] Metrics: %s | Dim: %s | Error: %v\n", metrics, dimensions, err)
		}
		return res, err
	}

	// Overview
	if res, err := fetchChannelMetrics(
		"views,estimatedMinutesWatched,averageViewDuration,averageViewPercentage", ""); err == nil && len(res.Rows) > 0 {
		row := res.Rows[0]
		analytics.Overview = OverviewStats{
			Views:          row[0].(float64),
			WatchTimeHours: row[1].(float64) / 60,
			AVD:            row[2].(float64),
			AVP:            row[3].(float64),
		}
	} else {
		return ChannelAnalyticsResponse{}, err
	}

	// Subscriber growth over time
	if res, err := fetchChannelMetrics(
		"subscribersGained,subscribersLost", "day"); err == nil {
		analytics.SubscriberGrowth = toRows(res)
	} else {
		return ChannelAnalyticsResponse{}, err
	}

	// Top videos by views
	topVideosLong, err := GetTopVideoIDs(startDate, endDate, 10, Longform.Ptr())
	if err != nil {
		return ChannelAnalyticsResponse{}, err
	}

	topVideosShort, err := GetTopVideoIDs(startDate, endDate, 10, Short.Ptr())
	if err != nil {
		return ChannelAnalyticsResponse{}, err
	}

	var topVideosWG sync.WaitGroup
	var outTopVideosLong = make([]VideoDetail, len(topVideosLong))
	var outTopVideosShort = make([]VideoDetail, len(topVideosShort))

	for idx, vidID := range topVideosLong {
		topVideosWG.Add(1)
		go func(i int, id string) {
			defer topVideosWG.Done()
			video, err := GetVideo(id)
			if err != nil {
				return
			}

			outTopVideosLong[i] = video
		}(idx, vidID)
	}

	for idx, vidID := range topVideosShort {
		topVideosWG.Add(1)
		go func(i int, id string) {
			defer topVideosWG.Done()
			video, err := GetVideo(id)
			if err != nil {
				return
			}

			outTopVideosShort[i] = video
		}(idx, vidID)
	}

	topVideosWG.Wait()
	analytics.TopVideosLong = outTopVideosLong
	analytics.TopVideosShort = outTopVideosShort

	// Traffic sources
	if res, err := fetchChannelMetrics(
		"views,estimatedMinutesWatched", "insightTrafficSourceType"); err == nil {
		analytics.TrafficSources = toRows(res)
	} else {
		return ChannelAnalyticsResponse{}, err
	}

	// Geography
	if res, err := fetchChannelMetrics(
		"views,estimatedMinutesWatched", "country"); err == nil {
		analytics.Geography = toRows(res)
	} else {
		return ChannelAnalyticsResponse{}, err
	}

	// Device types
	if res, err := fetchChannelMetrics(
		"views,estimatedMinutesWatched", "deviceType"); err == nil {
		analytics.DeviceTypes = toRows(res)
	} else {
		return ChannelAnalyticsResponse{}, err
	}

	// Age groups
	if res, err := fetchChannelMetrics(
		"viewerPercentage", "ageGroup"); err == nil {
		analytics.AgeGroups = toRows(res)
	} else {
		return ChannelAnalyticsResponse{}, err
	}

	// Gender
	if res, err := fetchChannelMetrics(
		"viewerPercentage", "gender"); err == nil {
		analytics.Gender = toRows(res)
	} else {
		return ChannelAnalyticsResponse{}, err
	}

	// Daily breakdown
	if res, err := fetchChannelMetrics(
		"views,estimatedMinutesWatched,subscribersGained,likes,shares", "day"); err == nil {
		analytics.DailyBreakdown = toRows(res)
	} else {
		return ChannelAnalyticsResponse{}, err
	}

	return analytics, nil
}
