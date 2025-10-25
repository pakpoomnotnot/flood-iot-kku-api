package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SensorData represents a flood sensor reading
type SensorData struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SensorID    string            `json:"sensor_id" bson:"sensor_id" binding:"required"`
	Location    Location          `json:"location" bson:"location" binding:"required"`
	WaterLevel  float64           `json:"water_level" bson:"water_level" binding:"required"`
	Temperature float64           `json:"temperature" bson:"temperature"`
	Humidity    float64           `json:"humidity" bson:"humidity"`
	Timestamp   time.Time         `json:"timestamp" bson:"timestamp"`
	Status      string            `json:"status" bson:"status"` // "normal", "warning", "critical"
	CreatedAt   time.Time         `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" bson:"updated_at"`
}

// Location represents geographical coordinates
type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" bson:"longitude" binding:"required"`
	Address   string  `json:"address,omitempty" bson:"address,omitempty"`
}

// SensorDataRequest represents the request payload for creating sensor data
type SensorDataRequest struct {
	SensorID    string    `json:"sensor_id" binding:"required"`
	Location    Location  `json:"location" binding:"required"`
	WaterLevel  float64   `json:"water_level" binding:"required"`
	Temperature  float64    `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Timestamp   time.Time `json:"timestamp"`
}

// SensorDataResponse represents the response payload for sensor data
type SensorDataResponse struct {
	ID          string    `json:"id"`
	SensorID    string    `json:"sensor_id"`
	Location    Location  `json:"location"`
	WaterLevel  float64   `json:"water_level"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	Timestamp   time.Time `json:"timestamp"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// PaginationResponse represents paginated response
type PaginationResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
}
