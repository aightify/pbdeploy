package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	internal "github.com/aightify/pbdeploy/internal/init"
)

var (
	version   = "v0.0.1"
	startTime = time.Now()
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		healthHandler(w, r)
	})

	router.HandleFunc("/deploy", internal.DeployHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Starting pbdeploy-agent on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Println("pbdeploy-agent started successfully")
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {

	w.Header().Add("Content-Type", "application/json")
	uptime := time.Since(startTime).String()

	healthz := map[string]string{
		"status":  "ok",
		"uptime":  uptime,
		"version": version,
	}

	json.NewEncoder(w).Encode(healthz)
}
