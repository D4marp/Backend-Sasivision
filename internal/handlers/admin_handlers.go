package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sasivision/backend/internal/models"
)

// slugify converts a title into a URL-friendly slug.
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			b.WriteRune(r)
			prevDash = false
		case r == ' ' || r == '-' || r == '_':
			if !prevDash {
				b.WriteRune('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

// --- File upload ---

var allowedUploadDirs = map[string]string{
	"markers":           "markers",
	"videos":            "videos",
	"audio":             "audio",
	"models":            "models",
	"videos/thumbnails": "videos/thumbnails",
}

func (h *Handler) UploadFile(c *gin.Context) {
	category := c.DefaultPostForm("category", "markers")
	subDir, ok := allowedUploadDirs[category]
	if !ok {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status: "error", Message: "Invalid upload category", Code: "ERR_INVALID_CATEGORY",
		})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status: "error", Message: "File is required", Code: "ERR_MISSING_FILE",
		})
		return
	}

	ext := filepath.Ext(file.Filename)
	base := slugify(strings.TrimSuffix(filepath.Base(file.Filename), ext))
	if base == "" {
		base = "file"
	}
	filename := fmt.Sprintf("%s-%d%s", base, time.Now().UnixNano(), strings.ToLower(ext))
	relPath := filepath.Join(subDir, filename)
	destPath := filepath.Join("storage", relPath)

	if err := c.SaveUploadedFile(file, destPath); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to save file", Code: "ERR_SAVE_FILE",
		})
		return
	}

	// Use forward slashes for URLs regardless of OS.
	urlPath := strings.ReplaceAll(relPath, "\\", "/")
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "File uploaded",
		Data: gin.H{"path": urlPath, "url": "/storage/" + urlPath},
	})
}

// --- Video CRUD ---

func (h *Handler) CreateVideo(c *gin.Context) {
	var req models.VideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c)
		return
	}
	if req.Slug == "" {
		req.Slug = slugify(req.Title)
	}
	id, err := h.content.CreateVideo(req)
	if err != nil {
		dbError(c, "Failed to create video")
		return
	}
	c.JSON(http.StatusCreated, models.ApiResponse{
		Status: "success", Message: "Video created", Data: gin.H{"id": id},
	})
}

func (h *Handler) UpdateVideo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	var req models.VideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c)
		return
	}
	if req.Slug == "" {
		req.Slug = slugify(req.Title)
	}
	if err := h.content.UpdateVideo(id, req); err != nil {
		dbError(c, "Failed to update video")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "Video updated"})
}

func (h *Handler) DeleteVideo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	if err := h.content.DeleteVideo(id); err != nil {
		dbError(c, "Failed to delete video")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "Video deleted"})
}

// --- Marker CRUD ---

func (h *Handler) CreateMarker(c *gin.Context) {
	var req models.MarkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c)
		return
	}
	if req.Slug == "" {
		req.Slug = slugify(req.Title)
	}
	id, err := h.content.CreateMarker(req)
	if err != nil {
		dbError(c, "Failed to create marker")
		return
	}
	c.JSON(http.StatusCreated, models.ApiResponse{
		Status: "success", Message: "Marker created", Data: gin.H{"id": id},
	})
}

func (h *Handler) UpdateMarker(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	var req models.MarkerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c)
		return
	}
	if req.Slug == "" {
		req.Slug = slugify(req.Title)
	}
	if err := h.content.UpdateMarker(id, req); err != nil {
		dbError(c, "Failed to update marker")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "Marker updated"})
}

func (h *Handler) DeleteMarker(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	if err := h.content.DeleteMarker(id); err != nil {
		dbError(c, "Failed to delete marker")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "Marker deleted"})
}

// --- Quiz category CRUD ---

func (h *Handler) GetAllQuizCategories(c *gin.Context) {
	categories, err := h.quiz.GetAllCategories()
	if err != nil {
		dbError(c, "Failed to load categories")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Categories retrieved", Data: categories,
	})
}

func (h *Handler) CreateQuizCategory(c *gin.Context) {
	var req models.QuizCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c)
		return
	}
	if req.Slug == "" {
		req.Slug = slugify(req.Name)
	}
	id, err := h.quiz.CreateCategory(req)
	if err != nil {
		dbError(c, "Failed to create category")
		return
	}
	c.JSON(http.StatusCreated, models.ApiResponse{
		Status: "success", Message: "Category created", Data: gin.H{"id": id},
	})
}

func (h *Handler) UpdateQuizCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	var req models.QuizCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c)
		return
	}
	if req.Slug == "" {
		req.Slug = slugify(req.Name)
	}
	if err := h.quiz.UpdateCategory(id, req); err != nil {
		dbError(c, "Failed to update category")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "Category updated"})
}

func (h *Handler) DeleteQuizCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	if err := h.quiz.DeleteCategory(id); err != nil {
		dbError(c, "Failed to delete category")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "Category deleted"})
}

// --- Quiz question CRUD ---

func (h *Handler) ListQuizQuestions(c *gin.Context) {
	categoryIDStr := c.Query("category_id")
	if categoryIDStr == "" {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status: "error", Message: "category_id query parameter is required", Code: "ERR_INVALID_REQUEST",
		})
		return
	}
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		badRequest(c)
		return
	}

	questions, err := h.quiz.GetQuestionsByCategoryID(categoryID)
	if err != nil {
		dbError(c, "Failed to load questions")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Questions retrieved", Data: questions,
	})
}

func (h *Handler) GetQuizQuestion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	question, err := h.quiz.GetQuestionByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status: "error", Message: "Question not found", Code: "ERR_NOT_FOUND",
		})
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Question retrieved", Data: question,
	})
}

func (h *Handler) CreateQuizQuestion(c *gin.Context) {
	var req models.QuizQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c)
		return
	}
	id, err := h.quiz.CreateQuestion(req)
	if err != nil {
		dbError(c, "Failed to create question")
		return
	}
	c.JSON(http.StatusCreated, models.ApiResponse{
		Status: "success", Message: "Question created", Data: gin.H{"id": id},
	})
}

func (h *Handler) UpdateQuizQuestion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	var req models.QuizQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c)
		return
	}
	if err := h.quiz.UpdateQuestion(id, req); err != nil {
		dbError(c, "Failed to update question")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "Question updated"})
}

func (h *Handler) DeleteQuizQuestion(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	if err := h.quiz.DeleteQuestion(id); err != nil {
		dbError(c, "Failed to delete question")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "Question deleted"})
}

// --- User management (admin only) ---

func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.users.List()
	if err != nil {
		dbError(c, "Failed to load users")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Users retrieved", Data: users,
	})
}

func (h *Handler) UpdateUserRole(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c)
		return
	}
	if err := h.users.UpdateRole(id, req.Role); err != nil {
		dbError(c, "Failed to update role")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "Role updated"})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		badRequest(c)
		return
	}
	if uid := currentUserID(c); uid != nil && *uid == id {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status: "error", Message: "You cannot delete your own account", Code: "ERR_SELF_DELETE",
		})
		return
	}
	if err := h.users.Delete(id); err != nil {
		dbError(c, "Failed to delete user")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "User deleted"})
}

// --- Analytics ---

// RecordAnalyticsEvent is a public endpoint clients call to log events.
func (h *Handler) RecordAnalyticsEvent(c *gin.Context) {
	var req models.AnalyticsEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		badRequest(c)
		return
	}
	userID := currentUserID(c)
	if err := h.analytics.RecordEvent(userID, req.EventType, req.EntityType, req.EntityID, req.Metadata); err != nil {
		dbError(c, "Failed to record event")
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{Status: "success", Message: "Event recorded"})
}

func (h *Handler) GetAnalytics(c *gin.Context) {
	now := time.Now()
	last7 := now.AddDate(0, 0, -7)
	last30 := now.AddDate(0, 0, -30)

	totalUsers, _ := h.users.Count()
	usersByRole, _ := h.users.CountByRole()
	newUsers7, _ := h.users.CountCreatedSince(last7)
	newUsers30, _ := h.users.CountCreatedSince(last30)

	totalQuestions, _ := h.quiz.CountQuestions()
	totalAttempts, _ := h.quiz.CountAttempts()
	attempts7, _ := h.quiz.CountAttemptsSince(last7)
	avgScore, _ := h.quiz.AverageScore()
	attemptsByCategory, _ := h.quiz.AttemptsByCategory()

	eventsByType, _ := h.analytics.CountByEventType()
	eventsDaily, _ := h.analytics.EventsDaily(14)
	topMarkers, _ := h.analytics.TopEntities("ar_scan", "markers", "title", 5)
	topVideos, _ := h.analytics.TopEntities("video_play", "videos", "title", 5)

	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Analytics retrieved",
		Data: gin.H{
			"users": gin.H{
				"total":      totalUsers,
				"by_role":    usersByRole,
				"new_last_7": newUsers7,
				"new_last_30": newUsers30,
			},
			"quiz": gin.H{
				"total_questions":      totalQuestions,
				"total_attempts":       totalAttempts,
				"attempts_last_7":      attempts7,
				"average_score":        avgScore,
				"attempts_by_category": attemptsByCategory,
			},
			"events": gin.H{
				"by_type":     eventsByType,
				"daily":       eventsDaily,
				"top_markers": topMarkers,
				"top_videos":  topVideos,
			},
		},
	})
}

// --- shared helpers ---

func badRequest(c *gin.Context) {
	c.JSON(http.StatusBadRequest, models.ApiResponse{
		Status: "error", Message: "Invalid request", Code: "ERR_INVALID_REQUEST",
	})
}

func dbError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, models.ApiResponse{
		Status: "error", Message: message, Code: "ERR_DB",
	})
}
