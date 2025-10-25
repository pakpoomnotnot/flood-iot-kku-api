package handlers

import (
	"flood-iot-kku-api/internal/database"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *database.MongoDB) {
	// Initialize repositories
	sensorRepo := database.NewSensorRepository(db, "sensor_data")

	// Initialize handlers
	sensorHandler := NewSensorHandler(sensorRepo)
	healthHandler := NewHealthHandler()

	// Health check routes
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/ready", healthHandler.ReadyCheck)

	// API routes
	api := router.Group("/api")
	{
		// Sensor data routes
		sensors := api.Group("/sensors")
		{
			sensors.POST("", sensorHandler.CreateSensorData)
			sensors.GET("", sensorHandler.GetAllSensorData)
			sensors.GET("/:id", sensorHandler.GetSensorData)
			sensors.PUT("/:id", sensorHandler.UpdateSensorData)
			sensors.DELETE("/:id", sensorHandler.DeleteSensorData)
			sensors.GET("/:sensor_id/latest", sensorHandler.GetLatestSensorData)
		}
	}
}
