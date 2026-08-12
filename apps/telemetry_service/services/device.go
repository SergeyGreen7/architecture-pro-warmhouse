package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	// "strings"
)

// DeviceService handles device service API
type DeviceService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// TemperatureResponse represents the response from the temperature API
type TemperatureSensor struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Status      string    `json:"status"`
	LastUpdated time.Time `json:"last_updated"`
	CreatedAt   time.Time `json:"created_at"`
}

// NewDeviceService creates a new temperature service
func NewDeviceService(baseURL string) *DeviceService {
	return &DeviceService{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *DeviceService) getTemperatureSensors(url string) ([]TemperatureSensor, error) {
	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error fetching sensors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var sensors []TemperatureSensor
	if err := json.NewDecoder(resp.Body).Decode(&sensors); err != nil {
		return nil, fmt.Errorf("error decoding sensors response: %w", err)
	}

	return sensors, nil
}

func (s *DeviceService) getTemperatureSensor(url string) (TemperatureSensor, error) {
	resp, err := s.HTTPClient.Get(url)
	if err != nil {
		return TemperatureSensor{}, fmt.Errorf("error fetching sensors: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return TemperatureSensor{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var sensor TemperatureSensor
	if err := json.NewDecoder(resp.Body).Decode(&sensor); err != nil {
		return TemperatureSensor{}, fmt.Errorf("error decoding sensors response: %w", err)
	}

	return sensor, nil
}

// Get temperature sensors
func (s *DeviceService) GetTemperatureSensors() ([]TemperatureSensor, error) {
	url := fmt.Sprintf("%s/api/v1/devices", s.BaseURL)
	log.Println("url =", url)

	return s.getTemperatureSensors(url)
}

// Get temperature sensors by ID
func (s *DeviceService) GetTemperatureSensorsByID(id int) (TemperatureSensor, error) {
	url := fmt.Sprintf("%s/api/v1/devices/id/%d", s.BaseURL, id)
	log.Println("url =", url)

	return s.getTemperatureSensor(url)
}

// Get temperature sensors by ID
func (s *DeviceService) GetTemperatureSensorsByLocation(location string) ([]TemperatureSensor, error) {
	url := fmt.Sprintf("%s/api/v1/devices/location/%s", s.BaseURL, location)
	log.Println("url =", url)

	return s.getTemperatureSensors(url)
}
