package youtube

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ChezyName/YouTube-MCP/config"
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
//TODO: use waitgroup to load all data async (all at once in threads)
func GetAnalyticsForVideo(videoID string, startDate string, endDate string, inRange string) (AnalyticsResponse, error) {
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
		return AnalyticsResponse{}, err
	}

	svc, err := youtubeAnalytics.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return AnalyticsResponse{}, err
	}

	analytics := AnalyticsResponse{
		VideoID:   videoID,
		DateRange: DateRange{Start: startDate, End: endDate},
	}

	var analyticsWaitgroup sync.WaitGroup
	var analyticsMutex sync.Mutex

	runTask := func(task func()) {
		analyticsWaitgroup.Add(1)
		go func() {
			defer analyticsWaitgroup.Done()
			task()
		}()
	}

	// Metrics for the main stats — views, watch time, AVD, AVP
	runTask(func() {
		if res, err := fetchMetrics(svc, videoID, startDate, endDate,
			"views,estimatedMinutesWatched,averageViewDuration,averageViewPercentage", ""); err == nil && len(res.Rows) > 0 {
			row := res.Rows[0]
			analyticsMutex.Lock()
			analytics.Overview = OverviewStats{
				Views:          row[0].(float64),
				WatchTimeHours: row[1].(float64) / 60,
				AVD:            row[2].(float64),
				AVP:            row[3].(float64),
			}
			analyticsMutex.Unlock()
		}
	})

	// Metrics for other stats (dislikes removed - must use returnYouTubeDislikes) — likes, dislikes, comments, shares, subs
	runTask(func() {
		if res, err := fetchMetrics(svc, videoID, startDate, endDate,
			"likes,dislikes,comments,shares,subscribersGained", ""); err == nil && len(res.Rows) > 0 {
			var trueDislikes = -1
			dislikes, err := fetchDislikes(videoID)
			if err == nil {
				trueDislikes = dislikes.Dislikes
			}

			row := res.Rows[0]
			analyticsMutex.Lock()
			analytics.Engagement = EngagementStats{
				Likes:       row[0].(float64),
				Dislikes:    float64(trueDislikes), //-1 for invalid / error
				Comments:    row[2].(float64),
				Shares:      row[3].(float64),
				Subscribers: row[4].(float64),
			}
			analyticsMutex.Unlock()
		}
	})

	// Impressions + CTR
	runTask(func() {
		if res, err := fetchMetrics(svc, videoID, startDate, endDate,
			"videoThumbnailImpressions,videoThumbnailImpressionsClickRate,uniqueViewers", ""); err == nil && len(res.Rows) > 0 {
			row := res.Rows[0]
			analyticsMutex.Lock()
			analytics.Impressions = ImpressionStats{
				Impressions: row[0].(float64),
				CTR:         row[1].(float64),
			}
			analyticsMutex.Unlock()
		}
	})

	// Traffic sources
	runTask(func() {
		if res, err := fetchMetrics(svc, videoID, startDate, endDate,
			"views,estimatedMinutesWatched", "insightTrafficSourceType"); err == nil {
			analyticsMutex.Lock()
			analytics.TrafficSources = toRows(res)
			analyticsMutex.Unlock()
		}
	})

	// Retention graph (full curve)
	runTask(func() {
		if res, err := fetchMetrics(svc, videoID, startDate, endDate,
			"audienceWatchRatio,relativeRetentionPerformance", "elapsedVideoTimeRatio"); err == nil {
			analyticsMutex.Lock()
			analytics.Retention = toRows(res)
			analyticsMutex.Unlock()
		}
	})

	// Geography
	runTask(func() {
		if res, err := fetchMetrics(svc, videoID, startDate, endDate,
			"views,estimatedMinutesWatched", "country"); err == nil {
			analyticsMutex.Lock()
			analytics.Geography = toRows(res)
			analyticsMutex.Unlock()
		}
	})

	// Device types
	runTask(func() {
		if res, err := fetchMetrics(svc, videoID, startDate, endDate,
			"views,estimatedMinutesWatched", "deviceType"); err == nil {
			analyticsMutex.Lock()
			analytics.DeviceTypes = toRows(res)
			analyticsMutex.Unlock()
		}
	})

	// Daily breakdown
	runTask(func() {
		if res, err := fetchMetrics(svc, videoID, startDate, endDate,
			"views,estimatedMinutesWatched,likes,shares", "day"); err == nil {
			analyticsMutex.Lock()
			analytics.DailyBreakdown = toRows(res)
			analyticsMutex.Unlock()
		}
	})

	analyticsWaitgroup.Wait()
	return analytics, nil
}

func GetTopVideoIDs(startDate string, endDate string, Limit int64) ([]string, error) {
	client, err := config.GetOAuthClient()
	if err != nil {
		return nil, err
	}

	svc, err := youtubeAnalytics.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	// This is the ONLY guaranteed-stable "Top Videos" query
	call := svc.Reports.Query().
		Ids("channel==MINE").
		StartDate(startDate).
		EndDate(endDate).
		Metrics("views").    // Only query views to find the "top" list
		Dimensions("video"). // Group by video
		Sort("-views").      // Sort descending
		MaxResults(Limit)    // Get top X

	res, err := call.Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[TOP VIDEOS ERROR]: %v\n", err)
		return nil, err
	}

	var videoIDs []string
	for _, row := range res.Rows {
		// row[0] is the video ID
		if id, ok := row[0].(string); ok {
			videoIDs = append(videoIDs, id)
		}
	}

	return videoIDs, nil
}
