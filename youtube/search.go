package youtube

import (
	"context"
	"fmt"
	"sync"

	"github.com/ChezyName/YouTube-MCP/config"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func SearchVideos(query string, limit int64, detailed bool, selfChannel bool) ([]*Video, []*VideoDetail, error) {
	client, err := config.GetOAuthClient()
	if err != nil {
		return nil, nil, err
	}
	ctx := context.Background()
	dataSvc, _ := youtube.NewService(ctx, option.WithHTTPClient(client))

	var call *youtube.SearchListCall
	if selfChannel {
		chanCall, err := dataSvc.Channels.List([]string{"id"}).
			ForHandle(config.GetConfig().ChannelHandle).
			Do()

		if err != nil || len(chanCall.Items) == 0 {
			return nil, nil, fmt.Errorf("failed to resolve handle to ID: %v", err)
		}

		// Need this to search for Self Videos
		actualChannelID := chanCall.Items[0].Id

		call = dataSvc.Search.List([]string{"snippet"}).
			ChannelId(actualChannelID).
			Q(query).
			Type("video").
			MaxResults(limit).
			SafeSearch("none")
	} else {
		call = dataSvc.Search.List([]string{"snippet"}).
			Q(query).
			Type("video").
			MaxResults(limit).
			SafeSearch("none")
	}

	res, err := call.Do()
	if err != nil {
		return nil, nil, err
	}

	var videoIDs []string

	var basicVideos = make([]*Video, len(res.Items))
	type result struct {
		index int
		vType VideoType
	}
	resultsChan := make(chan result)
	var typeWG sync.WaitGroup

	for i, item := range res.Items {
		videoIDs = append(videoIDs, item.Id.VideoId)
		if !detailed {
			//get the basics
			bVideo := Video{
				ID:          item.Id.VideoId,
				Title:       item.Snippet.Title,
				Description: item.Snippet.Description,
				PublishedAt: item.Snippet.PublishedAt,
				Thumbnail:   item.Snippet.Thumbnails.Medium.Url,
				Type:        Unknown,
			}
			basicVideos[i] = &bVideo
			typeWG.Add(1)
			go func(i int, vidID string) {
				defer typeWG.Done()

				videoType := Longform
				if isShort(vidID) {
					videoType = Short
				}

				resultsChan <- result{index: i, vType: videoType}
			}(i, item.Id.VideoId)
		}
	}

	if !detailed {
		go func() {
			typeWG.Wait()
			close(resultsChan)
		}()

		for res := range resultsChan {
			basicVideos[res.index].Type = res.vType
		}

		return basicVideos, nil, nil
	}

	var detailedVideos = make([]*VideoDetail, len(videoIDs))
	var wg sync.WaitGroup

	if detailed {
		for i, id := range videoIDs {
			wg.Add(1)
			go func(identifier string, index int) {
				defer wg.Done()
				details, err := GetVideo(identifier)
				if err != nil {
					detailedVideos[index] = nil
				} else {
					detailedVideos[index] = &details
				}

			}(id, i)
		}
	} else {

	}

	wg.Wait()
	return nil, detailedVideos, nil
}
