package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

const maxImageSize = 5 * 1024 * 1024 // 5 MB

type Status struct {
	Application string `json:"application"`
	Backend     string `json:"backend"`
	Database    string `json:"database"`
}

type Image struct {
	ID          int64     `json:"id"`
	Filename    string    `json:"filename"`
	ContentType string    `json:"content_type"`
	UploadedAt  time.Time `json:"uploaded_at"`
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

// Create the images table automatically.
func initializeDatabase() error {
	query := `
	CREATE TABLE IF NOT EXISTS images (
		id BIGSERIAL PRIMARY KEY,
		filename TEXT NOT NULL,
		content_type TEXT NOT NULL,
		data BYTEA NOT NULL,
		uploaded_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
	);`

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := db.ExecContext(ctx, query)
	return err
}

// Upload an image or list all uploaded images.
func imagesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		uploadImage(w, r)

	case http.MethodGet:
		listImages(w, r)

	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Upload image into PostgreSQL BYTEA.
func uploadImage(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxImageSize+1024*1024)

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Please upload an image using the image field", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(io.LimitReader(file, maxImageSize+1))
	if err != nil {
		http.Error(w, "Unable to read image", http.StatusBadRequest)
		return
	}

	if len(data) == 0 || len(data) > maxImageSize {
		http.Error(w, "Image must be between 1 byte and 5 MB", http.StatusBadRequest)
		return
	}

	contentType := http.DetectContentType(data)

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	if !allowedTypes[contentType] {
		http.Error(w, "Only JPEG, PNG, GIF and WebP images are allowed", http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO images (filename, content_type, data)
	VALUES ($1, $2, $3)
	RETURNING id`

	var imageID int64

	err = db.QueryRow(query, header.Filename, contentType, data).Scan(&imageID)
	if err != nil {
		log.Println("Image upload error:", err)
		http.Error(w, "Failed to save image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Image uploaded successfully",
		"id":      imageID,
	})
}

// Return image metadata without the image bytes.
func listImages(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		SELECT id, filename, content_type, uploaded_at
		FROM images
		ORDER BY uploaded_at DESC
	`)

	if err != nil {
		http.Error(w, "Unable to retrieve images", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	images := []Image{}

	for rows.Next() {
		var img Image

		err := rows.Scan(&img.ID, &img.Filename, &img.ContentType, &img.UploadedAt)
		if err != nil {
			http.Error(w, "Unable to read image records", http.StatusInternalServerError)
			return
		}

		images = append(images, img)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Unable to read image records", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(images)
}

// Retrieve an image or delete it using its ID.
func imageByIDHandler(w http.ResponseWriter, r *http.Request) {
	idText := strings.TrimPrefix(r.URL.Path, "/api/images/")

	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid image ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getImage(w, r, id)

	case http.MethodDelete:
		deleteImage(w, r, id)

	default:
		w.Header().Set("Allow", "GET, DELETE")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Retrieve image bytes from PostgreSQL.
func getImage(w http.ResponseWriter, r *http.Request, id int64) {
	var filename string
	var contentType string
	var data []byte

	err := db.QueryRow(`
		SELECT filename, content_type, data
		FROM images
		WHERE id = $1
	`, id).Scan(&filename, &contentType, &data)

	if err == sql.ErrNoRows {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Unable to retrieve image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	w.Write(data)
}

// Delete an image from PostgreSQL.
func deleteImage(w http.ResponseWriter, r *http.Request, id int64) {
	result, err := db.Exec("DELETE FROM images WHERE id = $1", id)

	if err != nil {
		http.Error(w, "Unable to delete image", http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Unable to confirm deletion", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Image deleted successfully",
	})
}

// Allow local React development and same-origin browser requests.
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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
		dbUser = "capstoneuser"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "capstonedb"
	}

	databaseURL := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPassword),
		Host:   net.JoinHostPort(dbHost, dbPort),
		Path:   dbName,
	}

	query := databaseURL.Query()
	query.Set("sslmode", "disable")
	databaseURL.RawQuery = query.Encode()

	var err error
	db, err = sql.Open("pgx", databaseURL.String())
	if err != nil {
		log.Fatal("Database configuration error:", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatal("Cannot connect to PostgreSQL:", err)
	}

	log.Println("Connected to PostgreSQL")

	if err := initializeDatabase(); err != nil {
		log.Fatal("Database initialization failed:", err)
	}

	log.Println("Images table is ready")

	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/ready", readyHandler)
	http.HandleFunc("/api/status", statusHandler)
	http.HandleFunc("/api/images", imagesHandler)
	http.HandleFunc("/api/images/", imageByIDHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	handler := corsMiddleware(http.DefaultServeMux)

	log.Printf("Backend server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
