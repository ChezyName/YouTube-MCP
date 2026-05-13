package router

import (
	"net/http"

	"github.com/ChezyName/YouTube-MCP/youtube"
	"github.com/gorilla/mux"
)

func CreateRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/health", health)
	router.HandleFunc("/videos", youtube.ListVideos)

	return router
}

func health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
