package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sasivision/backend/internal/database"
	"github.com/sasivision/backend/internal/models"
)

// Health reports API and database schema status.
func (h *Handler) Health(c *gin.Context) {
	ready, _ := database.SchemaReady(h.db)
	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"message":      "SasiVision API is running",
		"schema_ready": ready,
	})
}

// SetupStatus returns whether the database schema is initialized.
func (h *Handler) SetupStatus(c *gin.Context) {
	ready, err := database.SchemaReady(h.db)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ApiResponse{
			Status: "error", Message: "Database unavailable", Code: "ERR_DB",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":         "ok",
		"schema_ready":   ready,
		"migrations_dir": database.MigrationsDir(),
	})
}

// SetupInitialize runs migrations when schema is missing or when a valid setup token is provided.
func (h *Handler) SetupInitialize(c *gin.Context) {
	ready, err := database.SchemaReady(h.db)
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, models.ApiResponse{
			Status: "error", Message: "Database unavailable", Code: "ERR_DB",
		})
		return
	}

	token := c.GetHeader("X-Setup-Token")
	hasToken := h.cfg.SetupSecret != "" && token == h.cfg.SetupSecret

	if ready && !hasToken {
		c.JSON(http.StatusOK, models.ApiResponse{
			Status: "success", Message: "Database already initialized",
		})
		return
	}

	if err := database.RunMigrations(h.db, database.MigrationsDir()); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: err.Error(), Code: "ERR_MIGRATE",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Database initialized successfully",
	})
}
