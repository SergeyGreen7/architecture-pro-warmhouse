package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// TemperatureService fetches readings from the temperature API.
type TemperatureService struct {
	BaseURL    string
	HTTPClient *http.Client
}

// TemperatureResponse is the response returned by the temperature API.
type TemperatureResponse struct {
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	Timestamp   time.Time `json:"timestamp"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	SensorID    string    `json:"sensor_id"`
	SensorType  string    `json:"sensor_type"`
	Description string    `json:"description"`
}

func NewTemperatureService(baseURL string) *TemperatureService {
	return &TemperatureService{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *TemperatureService) GetTemperature(location string) (*TemperatureResponse, error) {
	requestURL, err := url.Parse(s.BaseURL + "/temperature")
	if err != nil {
		return nil, fmt.Errorf("invalid temperature API URL: %w", err)
	}

	query := requestURL.Query()
	query.Set("location", location)
	requestURL.RawQuery = query.Encode()

	return s.get(requestURL.String())
}

func (s *TemperatureService) GetTemperatureByID(sensorID string) (*TemperatureResponse, error) {
	return s.get(s.BaseURL + "/temperature/" + url.PathEscape(sensorID))
}

func (s *TemperatureService) get(requestURL string) (*TemperatureResponse, error) {
	resp, err := s.HTTPClient.Get(requestURL)
	if err != nil {
		return nil, fmt.Errorf("fetch temperature data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("temperature API returned status %d", resp.StatusCode)
	}

	var result TemperatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode temperature response: %w", err)
	}

	return &result, nil
}
