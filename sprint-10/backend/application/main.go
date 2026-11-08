package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Status struct {
	Application string `json:"application"`
	Backend     string `json:"backend"`
	Database    string `json:"database"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Backend is healthy"))
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Backend is ready"))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	response := Status{
		Application: "DevOps Capstone",
		Backend:     "Running",
		Database:    "Not connected yet",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ready", readyHandler)
	http.HandleFunc("/api/status", statusHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Backend server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
