package handlers

import (
	"net/http"

	"flood-iot-kku-api/internal/models"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles GET /health
func (h *HealthHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "API is healthy",
		Data: map[string]interface{}{
			"status":    "ok",
			"timestamp": "2024-01-01T00:00:00Z",
		},
	})
}

// ReadyCheck handles GET /ready
func (h *HealthHandler) ReadyCheck(c *gin.Context) {
	// Here you could add database connectivity checks
	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "API is ready",
		Data: map[string]interface{}{
			"status":    "ready",
			"timestamp": "2024-01-01T00:00:00Z",
		},
	})
}
