package youtube

import (
	"fmt"

	"github.com/rahadiangg/youtube-transcript-go/youtube"
)

func GetVideoTranscript(videoID string) (youtube.FetchedTranscript, error) {
	api := youtube.NewYouTubeTranscriptApi()

	// Fetch transcript (requesting English)
	transcript, err := api.Fetch(videoID, []string{"en"}, false)
	if err != nil {
		return youtube.FetchedTranscript{}, err
	}

	if transcript != nil {
		return *transcript, err
	}

	return youtube.FetchedTranscript{}, fmt.Errorf("YouTube Transcript API gave no error, but returned no transcript.")
}
