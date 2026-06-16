package server

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/sasivision/backend/internal/config"
	"github.com/sasivision/backend/internal/handlers"
	"github.com/sasivision/backend/internal/middleware"
)

// SetupRouter registers all SasiVision API routes on a Gin engine.
func SetupRouter(h *handlers.Handler, cfg *config.Config) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware(cfg))
	router.Use(middleware.RateLimitMiddleware(cfg))

	storagePath := filepath.Join(".", "storage")
	if _, err := os.Stat(storagePath); err == nil {
		router.Static("/storage", storagePath)
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "SasiVision API is running",
		})
	})

	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/sign-in", h.SignIn)
			auth.POST("/sign-up", h.SignUp)
			auth.POST("/verify-token", h.VerifyToken)
			auth.POST("/logout", middleware.AuthMiddleware(cfg), h.Logout)
		}

		quiz := api.Group("/quiz")
		{
			quiz.GET("/categories", h.GetQuizCategories)
			quiz.GET("/questions/:category", h.GetQuizQuestions)
			quiz.POST("/attempts", middleware.AuthMiddleware(cfg), h.SubmitQuizAttempt)
			quiz.GET("/history/:email", middleware.AuthMiddleware(cfg), h.GetQuizHistory)
			quiz.GET("/attempts/:id/details", middleware.AuthMiddleware(cfg), h.GetAttemptDetails)
		}

		api.POST("/analytics/event", middleware.OptionalAuthMiddleware(cfg), h.RecordAnalyticsEvent)

		content := api.Group("/content")
		{
			content.GET("/markers", h.GetMarkers)
			content.GET("/markers/:id", h.GetMarkerDetail)
			content.GET("/videos", h.GetVideos)
			content.GET("/videos/:id", h.GetVideoDetail)
		}

		features := api.Group("/features")
		{
			features.GET("/switches/:feature", h.GetFeatureSwitch)
			features.GET("/switches", h.GetAllFeatureSwitches)
		}

		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(cfg))

		manage := admin.Group("")
		manage.Use(middleware.RequireRole("admin", "editor"))
		{
			manage.POST("/upload", h.UploadFile)

			manage.POST("/videos", h.CreateVideo)
			manage.PUT("/videos/:id", h.UpdateVideo)
			manage.DELETE("/videos/:id", h.DeleteVideo)

			manage.POST("/markers", h.CreateMarker)
			manage.PUT("/markers/:id", h.UpdateMarker)
			manage.DELETE("/markers/:id", h.DeleteMarker)

			manage.GET("/quiz/categories", h.GetAllQuizCategories)
			manage.POST("/quiz/categories", h.CreateQuizCategory)
			manage.PUT("/quiz/categories/:id", h.UpdateQuizCategory)
			manage.DELETE("/quiz/categories/:id", h.DeleteQuizCategory)

			manage.GET("/quiz/questions", h.ListQuizQuestions)
			manage.GET("/quiz/questions/:id", h.GetQuizQuestion)
			manage.POST("/quiz/questions", h.CreateQuizQuestion)
			manage.PUT("/quiz/questions/:id", h.UpdateQuizQuestion)
			manage.DELETE("/quiz/questions/:id", h.DeleteQuizQuestion)
		}

		adminOnly := admin.Group("")
		adminOnly.Use(middleware.RequireRole("admin"))
		{
			adminOnly.PATCH("/features/:feature/enable", h.EnableFeature)
			adminOnly.PATCH("/features/:feature/disable", h.DisableFeature)
			adminOnly.GET("/stats", h.GetStats)
			adminOnly.GET("/analytics", h.GetAnalytics)

			adminOnly.GET("/users", h.ListUsers)
			adminOnly.PUT("/users/:id/role", h.UpdateUserRole)
			adminOnly.DELETE("/users/:id", h.DeleteUser)
		}
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Endpoint not found",
			"code":    "ERR_NOT_FOUND",
		})
	})

	return router
}
