package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLocationAndSensorIDMapping(t *testing.T) {
	tests := []struct {
		location string
		sensorID string
	}{
		{"Living Room", "1"},
		{"Bedroom", "2"},
		{"Kitchen", "3"},
		{"Unknown", "0"},
	}

	for _, test := range tests {
		if got := locationBySensorID(test.sensorID); got != test.location {
			t.Errorf("locationBySensorID(%q) = %q, want %q", test.sensorID, got, test.location)
		}
		if got := sensorIDByLocation(test.location); got != test.sensorID {
			t.Errorf("sensorIDByLocation(%q) = %q, want %q", test.location, got, test.sensorID)
		}
	}
}

func TestTemperatureHandler(t *testing.T) {
	tests := []string{
		"/temperature?location=Living%20Room",
		"/temperature/1",
	}

	for _, path := range tests {
		request := httptest.NewRequest(http.MethodGet, path, nil)
		response := httptest.NewRecorder()

		temperatureHandler(response, request)

		if response.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
		}

		var body temperatureResponse
		if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if body.Location != "Living Room" || body.SensorID != "1" {
			t.Fatalf("location = %q, sensor ID = %q", body.Location, body.SensorID)
		}
		if body.Value < 18 || body.Value >= 28 {
			t.Fatalf("temperature = %f, want value in [18, 28)", body.Value)
		}
	}
}
