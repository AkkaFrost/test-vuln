package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Secrets stored in variables
var (
	dbConnString = "postgres://admin:S3cur3P@ssw0rd!@prod-db.us-east-1.rds.amazonaws.com:5432/production?sslmode=disable"
	gcpAPIKey    = "AIzaSyD-Example1234567890abcdefGHIJKLMNO"
)

// --- Using the variables (2 times each) ---

func connectDB() *sql.DB {
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	return db
}

func reconnectDB() *sql.DB {
	log.Println("Reconnecting to database...")
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Reconnect failed: %v", err)
	}
	return db
}

func geocode(address string) (*http.Response, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", address, gcpAPIKey)
	return http.Get(url)
}

func getPlaceDetails(placeID string) (*http.Response, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/details/json?place_id=%s&key=%s", placeID, gcpAPIKey)
	return http.Get(url)
}

// --- Using the same values inline without the variables (2 times each) ---

func migrationJob() {
	db, err := sql.Open("postgres", "postgres://admin:S3cur3P@ssw0rd!@prod-db.us-east-1.rds.amazonaws.com:5432/production?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.Exec("ALTER TABLE users ADD COLUMN last_login TIMESTAMP")
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("postgres", "postgres://admin:S3cur3P@ssw0rd!@prod-db.us-east-1.rds.amazonaws.com:5432/production?sslmode=disable")
	if err != nil {
		http.Error(w, "DB down", 500)
		return
	}
	defer db.Close()
	fmt.Fprintf(w, "OK")
}

func reverseGeocode(lat, lng float64) (*http.Response, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?latlng=%f,%f&key=%s", lat, lng, "AIzaSyD-Example1234567890abcdefGHIJKLMNO")
	return http.Get(url)
}

func nearbySearch(lat, lng float64) (*http.Response, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%f,%f&radius=500&key=%s", lat, lng, "AIzaSyD-Example1234567890abcdefGHIJKLMNO")
	return http.Get(url)
}

func main() {
	db := connectDB()
	defer db.Close()

	http.HandleFunc("/health", healthCheck)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
