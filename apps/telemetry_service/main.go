package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"telemetry_service/handlers"
	"telemetry_service/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize temperature service
	temperatureAPIURL := getEnv("TEMPERATURE_API_URL", "http://temperature-api:8081")
	temperatureService := services.NewTemperatureService(temperatureAPIURL)
	log.Printf("Temperature service initialized with API URL: %s\n", temperatureAPIURL)

	// Initialize device service
	deviceAPIURL := getEnv("DEVICE_API_URL", "http://device-service:8083")
	deviceService := services.NewDeviceService(deviceAPIURL)
	log.Printf("Device service initialized with API URL: %s\n", deviceAPIURL)

	// Initialize router
	router := gin.Default()
	registerSwagger(router)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	apiRoutes := router.Group("/api/v1")

	// Register sensor routes
	telemetryHandler := handlers.NewTelemetryHandler(temperatureService, deviceService)
	telemetryHandler.RegisterRoutes(apiRoutes)

	// Start server
	srv := &http.Server{
		Addr:    getEnv("PORT", ":8082"),
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Telemetry service listening on port %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start Telemetry service: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down Telemetry service...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Service forced to shutdown: %v\n", err)
	}

	log.Println("Service exited properly")
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
