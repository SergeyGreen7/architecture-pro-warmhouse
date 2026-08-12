package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"telemetry_service/models"
	"telemetry_service/services"

	"github.com/gin-gonic/gin"
)

// TelemetryHandler handles sensor-related requests
type TelemetryHandler struct {
	TemperatureService *services.TemperatureService
	DeviceService      *services.DeviceService
}

// NewTelemetryHandler creates a new TelemetryHandler
func NewTelemetryHandler(
	temperatureService *services.TemperatureService,
	deviceService *services.DeviceService,
) *TelemetryHandler {
	return &TelemetryHandler{
		TemperatureService: temperatureService,
		DeviceService:      deviceService,
	}
}

// RegisterRoutes registers the sensor routes
func (h *TelemetryHandler) RegisterRoutes(router *gin.RouterGroup) {
	sensors := router.Group("/temperature")
	{
		sensors.GET("", h.GetTemperature)
		sensors.GET("/id/:id", h.GetTemperatureByID)
		sensors.GET("/location/:location", h.GetTemperatureByLocation)
	}
}

// GetTemperature handles GET /api/v1/temperature
func (h *TelemetryHandler) GetTemperature(c *gin.Context) {
	sensors, err := h.DeviceService.GetTemperatureSensors()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]models.TemperatureDscr, len(sensors))
	// Update temperature sensors with real-time data from the external API
	for i, sensor := range sensors {
		tempData, err := h.TemperatureService.GetTemperatureByID(fmt.Sprintf("%d", sensor.ID))
		if err == nil {
			result[i].ID = sensor.ID
			result[i].Unit = sensor.Unit
			result[i].Value = tempData.Value
			result[i].LastUpdated = tempData.Timestamp
			log.Printf("Updated temperature data for sensor %d from external API", sensor.ID)
		} else {
			log.Printf("Failed to fetch temperature data for sensor %d: %v", sensor.ID, err)
		}
	}

	c.JSON(http.StatusOK, result)
}

// GetTemperatureByID handles GET /api/v1/sensors/:id
func (h *TelemetryHandler) GetTemperatureByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sensor ID"})
		return
	}

	sensor, err := h.DeviceService.GetTemperatureSensorsByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sensor not found"})
		return
	}

	// If this is a temperature sensor, fetch real-time data from the temperature API
	var result models.TemperatureDscr
	tempData, err := h.TemperatureService.GetTemperatureByID(fmt.Sprintf("%d", sensor.ID))
	if err == nil {
		// Update sensor with real-time data
		result.ID = sensor.ID
		result.Unit = sensor.Unit
		result.Value = tempData.Value
		result.LastUpdated = tempData.Timestamp
		log.Printf("Updated temperature data for sensor %d from external API", sensor.ID)
	} else {
		log.Printf("Failed to fetch temperature data for sensor %d: %v", sensor.ID, err)
	}

	c.JSON(http.StatusOK, result)
}

// GetTemperatureByLocation handles GET /api/v1/sensors/temperature/:location
func (h *TelemetryHandler) GetTemperatureByLocation(c *gin.Context) {
	location := c.Param("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Location is required"})
		return
	}
	log.Println("location =", location)

	sensors, err := h.DeviceService.GetTemperatureSensorsByLocation(location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := make([]models.TemperatureDscr, len(sensors))
	// Update temperature sensors with real-time data from the external API
	for i, sensor := range sensors {
		tempData, err := h.TemperatureService.GetTemperatureByID(fmt.Sprintf("%d", sensor.ID))
		if err == nil {
			result[i].ID = sensor.ID
			result[i].Unit = sensor.Unit
			result[i].Value = tempData.Value
			result[i].LastUpdated = tempData.Timestamp
			log.Printf("Updated temperature data for sensor %d from external API", sensor.ID)
		} else {
			log.Printf("Failed to fetch temperature data for sensor %d: %v", sensor.ID, err)
		}
	}

	c.JSON(http.StatusOK, result)
}
