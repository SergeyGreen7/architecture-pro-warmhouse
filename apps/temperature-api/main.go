package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type temperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})
	http.HandleFunc("/temperature", temperatureHandler)
	http.HandleFunc("/temperature/", temperatureHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("temperature-api is listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func temperatureHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	location := r.URL.Query().Get("location")
	sensorID := r.URL.Query().Get("sensorId")
	if sensorID == "" {
		sensorID = r.URL.Query().Get("sensorID")
	}
	if sensorID == "" && strings.HasPrefix(r.URL.Path, "/temperature/") {
		sensorID = strings.TrimPrefix(r.URL.Path, "/temperature/")
	}

	if location == "" {
		location = locationBySensorID(sensorID)
	}
	if sensorID == "" {
		sensorID = sensorIDByLocation(location)
	}

	writeJSON(w, temperatureResponse{
		Value:       20 + rand.Float64()*10,
		Unit:        "°C",
		Timestamp:   time.Now().UTC(),
		Location:    location,
		Status:      "active",
		SensorID:    sensorID,
		SensorType:  "temperature",
		Description: "Temperature sensor in " + location,
	})
}

func locationBySensorID(sensorID string) string {
	switch sensorID {
	case "1":
		return "Living Room"
	case "2":
		return "Bedroom"
	case "3":
		return "Kitchen"
	default:
		return "Unknown"
	}
}

func sensorIDByLocation(location string) string {
	switch location {
	case "Living Room":
		return "1"
	case "Bedroom":
		return "2"
	case "Kitchen":
		return "3"
	default:
		return "0"
	}
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("could not encode response: %v", err)
	}
}
