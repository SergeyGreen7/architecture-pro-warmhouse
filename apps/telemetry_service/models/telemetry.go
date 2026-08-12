package models

import (
	"time"
)

// SensorType represents the type of sensor
type SensorType string

const (
	Temperature SensorType = "temperature"
)

// Sensor represents a smart home sensor
type Sensor struct {
	ID   int        `json:"id"`
	Type SensorType `json:"type"`
	Unit string     `json:"unit"`
}

// Sensor represents a smart home sensor
type TemperatureDscr struct {
	ID          int       `json:"id"`
	Value       float64   `json:"value"`
	Unit        string    `json:"unit"`
	LastUpdated time.Time `json:"last_updated"`
}
