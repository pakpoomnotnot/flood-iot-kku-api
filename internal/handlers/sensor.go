package handlers

import (
	"net/http"
	"strconv"
	"time"

	"flood-iot-kku-api/internal/database"
	"flood-iot-kku-api/internal/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SensorHandler struct {
	repo *database.SensorRepository
}

func NewSensorHandler(repo *database.SensorRepository) *SensorHandler {
	return &SensorHandler{
		repo: repo,
	}
}

// CreateSensorData handles POST /api/sensors
func (h *SensorHandler) CreateSensorData(c *gin.Context) {
	var req models.SensorDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	// Convert request to sensor data model
	sensorData := &models.SensorData{
		ID:          primitive.NewObjectID(),
		SensorID:    req.SensorID,
		Location:    req.Location,
		WaterLevel:  req.WaterLevel,
		Temperature: req.Temperature,
		Humidity:    req.Humidity,
		Timestamp:   req.Timestamp,
	}

	// If timestamp is not provided, use current time
	if sensorData.Timestamp.IsZero() {
		sensorData.Timestamp = time.Now()
	}

	if err := h.repo.Create(c.Request.Context(), sensorData); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create sensor data",
			Error:   err.Error(),
		})
		return
	}

	// Convert to response format
	response := h.convertToResponse(sensorData)

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Sensor data created successfully",
		Data:    response,
	})
}

// GetSensorData handles GET /api/sensors/:id
func (h *SensorHandler) GetSensorData(c *gin.Context) {
	id := c.Param("id")

	sensorData, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Sensor data not found",
			Error:   err.Error(),
		})
		return
	}

	response := h.convertToResponse(sensorData)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Sensor data retrieved successfully",
		Data:    response,
	})
}

// GetAllSensorData handles GET /api/sensors
func (h *SensorHandler) GetAllSensorData(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	status := c.Query("status")
	sensorID := c.Query("sensor_id")

	var results []models.SensorData
	var total int64
	var err error

	if status != "" {
		results, total, err = h.repo.GetByStatus(c.Request.Context(), status, page, limit)
	} else if sensorID != "" {
		results, total, err = h.repo.GetBySensorID(c.Request.Context(), sensorID, page, limit)
	} else {
		results, total, err = h.repo.GetAll(c.Request.Context(), page, limit)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to retrieve sensor data",
			Error:   err.Error(),
		})
		return
	}

	// Convert to response format
	var responseData []models.SensorDataResponse
	for _, data := range results {
		responseData = append(responseData, h.convertToResponse(&data))
	}

	totalPages := (total + limit - 1) / limit

	paginationResponse := models.PaginationResponse{
		Data:       responseData,
		Page:       int(page),
		Limit:      int(limit),
		Total:      total,
		TotalPages: int(totalPages),
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Sensor data retrieved successfully",
		Data:    paginationResponse,
	})
}

// UpdateSensorData handles PUT /api/sensors/:id
func (h *SensorHandler) UpdateSensorData(c *gin.Context) {
	id := c.Param("id")

	var req models.SensorDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request data",
			Error:   err.Error(),
		})
		return
	}

	// Get existing data first
	existingData, err := h.repo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Sensor data not found",
			Error:   err.Error(),
		})
		return
	}

	// Update fields
	existingData.SensorID = req.SensorID
	existingData.Location = req.Location
	existingData.WaterLevel = req.WaterLevel
	existingData.Temperature = req.Temperature
	existingData.Humidity = req.Humidity
	if !req.Timestamp.IsZero() {
		existingData.Timestamp = req.Timestamp
	}

	if err := h.repo.Update(c.Request.Context(), id, existingData); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to update sensor data",
			Error:   err.Error(),
		})
		return
	}

	response := h.convertToResponse(existingData)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Sensor data updated successfully",
		Data:    response,
	})
}

// DeleteSensorData handles DELETE /api/sensors/:id
func (h *SensorHandler) DeleteSensorData(c *gin.Context) {
	id := c.Param("id")

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete sensor data",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Sensor data deleted successfully",
	})
}

// GetLatestSensorData handles GET /api/sensors/:sensor_id/latest
func (h *SensorHandler) GetLatestSensorData(c *gin.Context) {
	sensorID := c.Param("sensor_id")

	sensorData, err := h.repo.GetLatestBySensorID(c.Request.Context(), sensorID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Latest sensor data not found",
			Error:   err.Error(),
		})
		return
	}

	response := h.convertToResponse(sensorData)

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Latest sensor data retrieved successfully",
		Data:    response,
	})
}

// convertToResponse converts SensorData to SensorDataResponse
func (h *SensorHandler) convertToResponse(data *models.SensorData) models.SensorDataResponse {
	return models.SensorDataResponse{
		ID:          data.ID.Hex(),
		SensorID:    data.SensorID,
		Location:    data.Location,
		WaterLevel:  data.WaterLevel,
		Temperature: data.Temperature,
		Humidity:    data.Humidity,
		Timestamp:   data.Timestamp,
		Status:      data.Status,
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}
}
