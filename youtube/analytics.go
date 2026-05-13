package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ChezyName/YouTube-MCP/config"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
	youtubeAnalytics "google.golang.org/api/youtubeanalytics/v2"
)

func fetchMetrics(svc *youtubeAnalytics.Service, videoID, startDate, endDate, metrics, dimensions string) (*youtubeAnalytics.QueryResponse, error) {
	call := svc.Reports.Query().
		Ids("channel==MINE").
		StartDate(startDate).
		EndDate(endDate).
		Metrics(metrics).
		Filters("video==" + videoID)

	if dimensions != "" {
		call = call.Dimensions(dimensions)
	}

	return call.Do()
}

// turns graph data into rows
func toRows(res *youtubeAnalytics.QueryResponse) []RowData {
	if res == nil {
		return nil
	}
	var rows []RowData
	for _, row := range res.Rows {
		r := RowData{}
		if len(row) > 0 {
			if s, ok := row[0].(string); ok {
				r.Label = s
			}
		}
		if len(row) > 1 {
			if v, ok := row[1].(float64); ok {
				r.Value = v
			}
		}
		if len(row) > 2 {
			if v, ok := row[2].(float64); ok {
				r.Value2 = v
			}
		}
		rows = append(rows, r)
	}
	return rows
}

/*
(defaults to 90days)
Params:

	start: start date for range, needs end date (defaults to -90days)
	end: end date for range
	range: custom numbers, in days or Lifetime - this superseeds all
*/
func GetAnalyticsForVideo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	videoID := vars["id"]

	//given URL...?start=XYZ&end=XYZ
	endDate := r.URL.Query().Get("end")
	startDate := r.URL.Query().Get("start")
	inRange := r.URL.Query().Get("range")

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

	analytics := AnalyticsResponse{
		VideoID:   videoID,
		DateRange: DateRange{Start: startDate, End: endDate},
	}

	// Metrics for the main stats — views, watch time, AVD, AVP
	if res, err := fetchMetrics(svc, videoID, startDate, endDate,
		"views,estimatedMinutesWatched,averageViewDuration,averageViewPercentage", ""); err == nil && len(res.Rows) > 0 {
		row := res.Rows[0]
		analytics.Overview = OverviewStats{
			Views:          row[0].(float64),
			WatchTimeHours: row[1].(float64) / 60,
			AVD:            row[2].(float64),
			AVP:            row[3].(float64),
		}
	}

	// Metrics for other stats (dislikes removed - must use returnYouTubeDislikes) — likes, dislikes, comments, shares, subs
	if res, err := fetchMetrics(svc, videoID, startDate, endDate,
		"likes,dislikes,comments,shares,subscribersGained", ""); err == nil && len(res.Rows) > 0 {
		var trueDislikes = -1
		dislikes, err := fetchDislikes(videoID)
		if err != nil {
			trueDislikes = dislikes.Dislikes
		}

		row := res.Rows[0]
		analytics.Engagement = EngagementStats{
			Likes:       row[0].(float64),
			Dislikes:    float64(trueDislikes), //-1 for invalid / error
			Comments:    row[2].(float64),
			Shares:      row[3].(float64),
			Subscribers: row[4].(float64),
		}
	}

	// Impressions + CTR
	if res, err := fetchMetrics(svc, videoID, startDate, endDate,
		"impressions,impressionClickThroughRate,uniqueViewers", ""); err == nil && len(res.Rows) > 0 {
		row := res.Rows[0]
		analytics.Impressions = ImpressionStats{
			Impressions: row[0].(float64),
			CTR:         row[1].(float64),
			UniqueViews: row[2].(float64),
		}
	}

	// Traffic sources
	if res, err := fetchMetrics(svc, videoID, startDate, endDate,
		"views,estimatedMinutesWatched", "insightTrafficSourceType"); err == nil {
		analytics.TrafficSources = toRows(res)
	}

	// Retention graph (full curve)
	if res, err := fetchMetrics(svc, videoID, startDate, endDate,
		"audienceWatchRatio,relativeRetentionPerformance", "elapsedVideoTimeRatio"); err == nil {
		analytics.Retention = toRows(res)
	}

	// Geography
	if res, err := fetchMetrics(svc, videoID, startDate, endDate,
		"views,estimatedMinutesWatched", "country"); err == nil {
		analytics.Geography = toRows(res)
	}

	// Device types
	if res, err := fetchMetrics(svc, videoID, startDate, endDate,
		"views,estimatedMinutesWatched", "deviceType"); err == nil {
		analytics.DeviceTypes = toRows(res)
	}

	// Daily breakdown
	if res, err := fetchMetrics(svc, videoID, startDate, endDate,
		"views,estimatedMinutesWatched,likes,shares", "day"); err == nil {
		analytics.DailyBreakdown = toRows(res)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analytics)
}
