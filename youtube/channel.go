package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ChezyName/YouTube-MCP/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
	youtubeAnalytics "google.golang.org/api/youtubeanalytics/v2"
)

func GetChannel(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	svc, err := youtube.NewService(ctx, option.WithAPIKey(config.GetConfig().YouTubeAPI))
	if err != nil {
		http.Error(w, "Failed to create YouTube client", http.StatusInternalServerError)
		return
	}

	res, err := svc.Channels.List([]string{"snippet", "statistics", "brandingSettings"}).
		ForHandle(config.GetConfig().ChannelHandle).
		Do()
	if err != nil {
		http.Error(w, "Failed to fetch channel: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(res.Items) == 0 {
		http.Error(w, "Channel not found", http.StatusNotFound)
		return
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
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(channel)
}

/*
(defaults to lifetime)
Params:

	start: start date for range, needs end date (defaults to -90days)
	end: end date for range
	range: custom numbers, in days or Lifetime - this superseeds all
*/
func GetChannelAnalytics(w http.ResponseWriter, r *http.Request) {
	//given URL...?start=XYZ&end=XYZ
	endDate := r.URL.Query().Get("end")
	startDate := r.URL.Query().Get("start")
	inRange := r.URL.Query().Get("range")

	//defaults to lifetiem
	if startDate == "" && endDate == "" && inRange == "" {
		endDate = time.Now().Format("2006-01-02")
		startDate = "2005-01-01"
	} else {
		if endDate == "" {
			endDate = time.Now().Format("2006-01-02")
		}
		if startDate == "" {
			//TODO: allow edits via .env or some kind of config.json to edit this
			startDate = time.Now().AddDate(0, 0, -90).Format("2006-01-02") //last 90
		}

		if inRange != "" {
			if strings.ToUpper(inRange) == "LIFETIME" {
				endDate = time.Now().Format("2006-01-02")
				startDate = "2005-01-01"
			} else {
				//try and parse # of days
				days := 90 // default
				fmt.Sscanf(inRange, "%d", &days)
				startDate = time.Now().AddDate(0, 0, -days).Format("2006-01-02")
				endDate = time.Now().Format("2006-01-02")
			}
		}
	}

	client, err := config.GetOAuthClient()
	if err != nil {
		http.Error(w, "OAuth client error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	svc, err := youtubeAnalytics.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		http.Error(w, "Failed to create Analytics client: "+err.Error(), http.StatusInternalServerError)
		return
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
		return call.Do()
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
	}

	// Impressions + CTR
	if res, err := fetchChannelMetrics(
		"impressions,impressionClickThroughRate,uniqueViewers", ""); err == nil && len(res.Rows) > 0 {
		row := res.Rows[0]
		analytics.Impressions = ImpressionStats{
			Impressions: row[0].(float64),
			CTR:         row[1].(float64),
			UniqueViews: row[2].(float64),
		}
	}

	// Subscriber growth over time
	if res, err := fetchChannelMetrics(
		"subscribersGained,subscribersLost", "day"); err == nil {
		analytics.SubscriberGrowth = toRows(res)
	}

	// Top videos by views
	if res, err := fetchChannelMetrics(
		"views,estimatedMinutesWatched,averageViewPercentage", "video"); err == nil {
		call := svc.Reports.Query().
			Ids("channel==MINE").
			StartDate(startDate).
			EndDate(endDate).
			Metrics("views,estimatedMinutesWatched,averageViewPercentage").
			Dimensions("video").
			Sort("-views").
			MaxResults(20)
		if topRes, err := call.Do(); err == nil {
			analytics.TopVideos = toRows(topRes)
		}
		_ = res
	}

	// Traffic sources
	if res, err := fetchChannelMetrics(
		"views,estimatedMinutesWatched", "insightTrafficSourceType"); err == nil {
		analytics.TrafficSources = toRows(res)
	}

	// Geography
	if res, err := fetchChannelMetrics(
		"views,estimatedMinutesWatched", "country"); err == nil {
		analytics.Geography = toRows(res)
	}

	// Device types
	if res, err := fetchChannelMetrics(
		"views,estimatedMinutesWatched", "deviceType"); err == nil {
		analytics.DeviceTypes = toRows(res)
	}

	// Age groups
	if res, err := fetchChannelMetrics(
		"viewerPercentage", "ageGroup"); err == nil {
		analytics.AgeGroups = toRows(res)
	}

	// Gender
	if res, err := fetchChannelMetrics(
		"viewerPercentage", "gender"); err == nil {
		analytics.Gender = toRows(res)
	}

	// Daily breakdown
	if res, err := fetchChannelMetrics(
		"views,estimatedMinutesWatched,subscribersGained,likes,shares", "day"); err == nil {
		analytics.DailyBreakdown = toRows(res)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}
