package main

import (
	"config"
	"handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	err := config.LoadConfig()

	if err != nil {
		log.Printf("error loading config: %v\n", err)
		os.Exit(1)
	}

  config.WatchConfig()

	http.HandleFunc("/robots.txt", handlers.RobotsTxtHandler)
	http.HandleFunc("/status", handlers.StatusHandler)
	http.HandleFunc("/", handlers.IndexHandler)
	err = http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Println("error starting listening on port :8080")
		os.Exit(1)
	}
}

