package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

type Status struct {
	Application string `json:"application"`
	Backend     string `json:"backend"`
	Database    string `json:"database"`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Backend is healthy"))
}

func checkDatabase() bool {
	if db == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return db.PingContext(ctx) == nil
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	if !checkDatabase() {
		http.Error(w, "Database is not ready", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Backend and database are ready"))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	databaseStatus := "Disconnected"

	if checkDatabase() {
		databaseStatus = "Connected"
	}

	response := Status{
		Application: "DevOps Capstone",
		Backend:     "Running",
		Database:    databaseStatus,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "appuser"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "appdb"
	}

	dsn := "host=" + dbHost +
		" port=" + dbPort +
		" user=" + dbUser +
		" password=" + dbPassword +
		" dbname=" + dbName +
		" sslmode=disable"

	var err error
	db, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Database configuration error: ", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

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
