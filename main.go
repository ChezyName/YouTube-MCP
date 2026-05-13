package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ChezyName/YouTube-MCP/config"
	"github.com/ChezyName/YouTube-MCP/router"
)

func init() {
	config.LoadConfig()
	//config.EnsureRefreshToken() - only allowed for clients
}

func main() {
	router := router.CreateRouter()

	server := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("Starting Server")
	log.Fatal(server.ListenAndServe())
}
