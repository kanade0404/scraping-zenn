package main

import (
	"github.com/kanade0404/scraping-zenn/logger"
	"net/http"
	"os"

	scraping_zenn "github.com/kanade0404/scraping-zenn"
)

func main() {
	port := "8888"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	http.HandleFunc("/scraping", scraping_zenn.Scraping)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		logger.Fatalf("failed to starting server. error: %v", err.Error())
	}
}
