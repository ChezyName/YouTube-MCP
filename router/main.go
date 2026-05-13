package router

import (
	"net/http"

	"github.com/ChezyName/YouTube-MCP/youtube"
	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", health).Methods("GET") //Health checker

	//Routes for Videos
	router.HandleFunc("/videos", youtube.ListVideos).Methods("GET")                          // get all videos from Channel Handle in env
	router.HandleFunc("/videos/{id}", youtube.GetVideo).Methods("GET")                       // get the specific video data such as title, desc, view count (from PUBLIC data)
	router.HandleFunc("/videos/{id}/analytics", youtube.GetAnalyticsForVideo).Methods("GET") // gets the specific video data from youtube analytics - needs OAUTH

	router.HandleFunc("/channel", youtube.GetChannel).Methods("GET")                    // gets the channel stats - the public data
	router.HandleFunc("/channel/analytics", youtube.GetChannelAnalytics).Methods("GET") // gets the channel stats - the public data

	return router
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
