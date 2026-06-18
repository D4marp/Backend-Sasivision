package handlers

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sasivision/backend/internal/config"
	"github.com/sasivision/backend/internal/models"
	"github.com/sasivision/backend/internal/repositories"
	"github.com/sasivision/backend/internal/utils"
)

type Handler struct {
	db        *sql.DB
	cfg       *config.Config
	users     *repositories.UserRepository
	sessions  *repositories.SessionRepository
	quiz      *repositories.QuizRepository
	content   *repositories.ContentRepository
	features  *repositories.FeatureRepository
	analytics *repositories.AnalyticsRepository
}

func New(db *sql.DB, cfg *config.Config) *Handler {
	return &Handler{
		db:        db,
		cfg:       cfg,
		users:     repositories.NewUserRepository(db),
		sessions:  repositories.NewSessionRepository(db),
		quiz:      repositories.NewQuizRepository(db),
		content:   repositories.NewContentRepository(db),
		features:  repositories.NewFeatureRepository(db),
		analytics: repositories.NewAnalyticsRepository(db),
	}
}

func (h *Handler) SignIn(c *gin.Context) {
	var req models.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status: "error", Message: "Invalid request", Code: "ERR_INVALID_REQUEST",
		})
		return
	}

	user, passwordHash, err := h.users.FindByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status: "error", Message: "Email or password is incorrect", Code: "ERR_INVALID_CREDENTIALS",
		})
		return
	}

	if !utils.CheckPassword(req.Password, passwordHash) {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status: "error", Message: "Email or password is incorrect", Code: "ERR_INVALID_CREDENTIALS",
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, h.cfg.JWTSecret, h.cfg.JWTExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to generate token", Code: "ERR_TOKEN_GENERATION",
		})
		return
	}

	if err := h.sessions.Create(user.ID, token, time.Now().Add(h.cfg.JWTExpiry)); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to create session", Code: "ERR_SESSION_CREATE",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Login successful",
		Data: gin.H{"email": user.Email, "token": token, "role": user.Role},
	})
}

func (h *Handler) SignUp(c *gin.Context) {
	var req models.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status: "error", Message: "Invalid request", Code: "ERR_INVALID_REQUEST",
		})
		return
	}

	if _, _, err := h.users.FindByEmail(req.Email); err == nil {
		c.JSON(http.StatusConflict, models.ApiResponse{
			Status: "error", Message: "Email already registered", Code: "ERR_EMAIL_EXISTS",
		})
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to hash password", Code: "ERR_HASH_PASSWORD",
		})
		return
	}

	user, err := h.users.Create(req.Email, hash, req.FullName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to create user", Code: "ERR_CREATE_USER",
		})
		return
	}

	c.JSON(http.StatusCreated, models.ApiResponse{
		Status: "success", Message: "Registration successful",
		Data: gin.H{"email": user.Email},
	})
}

func (h *Handler) VerifyToken(c *gin.Context) {
	token := extractToken(c)
	if token == "" {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status: "error", Message: "Token required", Code: "ERR_MISSING_TOKEN",
		})
		return
	}

	claims, err := utils.ParseToken(token, h.cfg.JWTSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status: "error", Message: "Invalid token", Code: "ERR_INVALID_TOKEN",
		})
		return
	}

	if _, err := h.sessions.FindValid(token); err != nil {
		c.JSON(http.StatusUnauthorized, models.ApiResponse{
			Status: "error", Message: "Session expired", Code: "ERR_SESSION_EXPIRED",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Token is valid",
		Data: gin.H{"valid": true, "email": claims.Email, "role": claims.Role},
	})
}

func (h *Handler) Logout(c *gin.Context) {
	token := extractToken(c)
	if token != "" {
		_ = h.sessions.DeleteByToken(token)
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Logged out successfully",
	})
}

func (h *Handler) GetQuizCategories(c *gin.Context) {
	categories, err := h.quiz.GetActiveCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to load categories", Code: "ERR_DB_QUERY",
		})
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Categories retrieved", Data: categories,
	})
}

func (h *Handler) GetQuizQuestions(c *gin.Context) {
	category := c.Param("category")
	questions, err := h.quiz.GetQuestionsByCategory(category)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status: "error", Message: "Category not found", Code: "ERR_NOT_FOUND",
		})
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Questions retrieved", Data: questions,
	})
}

func (h *Handler) SubmitQuizAttempt(c *gin.Context) {
	var req models.QuizSubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status: "error", Message: "Invalid request", Code: "ERR_INVALID_REQUEST",
		})
		return
	}

	// Ensure the authenticated user can only submit attempts for themselves.
	authEmail := c.GetString("email")
	if authEmail != "" && req.Email != authEmail {
		c.JSON(http.StatusForbidden, models.ApiResponse{
			Status: "error", Message: "Cannot submit attempt for another user", Code: "ERR_FORBIDDEN",
		})
		return
	}

	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)
	if startTime.IsZero() {
		startTime = time.Now()
	}
	if endTime.IsZero() {
		endTime = time.Now()
	}

	attempt := models.QuizAttempt{
		Email:        req.Email,
		CategoryID:   req.CategoryID,
		CorrectCount: req.Correct,
		TotalCount:   req.Total,
		Score:        req.Score,
		StartTime:    startTime,
		EndTime:      endTime,
		FinishDate:   req.FinishDate,
	}

	attemptID, err := h.quiz.CreateAttempt(attempt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to save attempt", Code: "ERR_DB_INSERT",
		})
		return
	}

	details := make([]models.AttemptDetail, 0, len(req.Answers))
	for _, answer := range req.Answers {
		details = append(details, models.AttemptDetail{
			QuizID:     answer.QuizID,
			Type:       answer.Type,
			UserAnswer: answer.Answers,
		})
	}
	if err := h.quiz.CreateAttemptDetails(attemptID, details); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to save attempt details", Code: "ERR_DB_INSERT",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Quiz submitted successfully",
		Data: gin.H{"attempt_id": attemptID, "score": req.Score},
	})
}

func (h *Handler) GetQuizHistory(c *gin.Context) {
	email := c.Param("email")

	// Users may only read their own history; admins/editors can read anyone's.
	authEmail := c.GetString("email")
	role := c.GetString("role")
	if authEmail != email && role != "admin" && role != "editor" {
		c.JSON(http.StatusForbidden, models.ApiResponse{
			Status: "error", Message: "Cannot view another user's history", Code: "ERR_FORBIDDEN",
		})
		return
	}

	attempts, err := h.quiz.GetAttemptsByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to load history", Code: "ERR_DB_QUERY",
		})
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "History retrieved", Data: attempts,
	})
}

func (h *Handler) GetAttemptDetails(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status: "error", Message: "Invalid attempt id", Code: "ERR_INVALID_REQUEST",
		})
		return
	}

	ownerEmail, err := h.quiz.GetAttemptOwnerEmail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status: "error", Message: "Attempt not found", Code: "ERR_NOT_FOUND",
		})
		return
	}

	authEmail := c.GetString("email")
	role := c.GetString("role")
	if authEmail != ownerEmail && role != "admin" && role != "editor" {
		c.JSON(http.StatusForbidden, models.ApiResponse{
			Status: "error", Message: "Cannot view another user's attempt details", Code: "ERR_FORBIDDEN",
		})
		return
	}

	details, err := h.quiz.GetAttemptDetails(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to load details", Code: "ERR_DB_QUERY",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Details retrieved",
		Data: gin.H{"attempt_id": id, "details": details},
	})
}

func (h *Handler) GetMarkers(c *gin.Context) {
	markers, err := h.content.GetMarkers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to load markers", Code: "ERR_DB_QUERY",
		})
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Markers retrieved", Data: markers,
	})
}

func (h *Handler) GetMarkerDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status: "error", Message: "Invalid marker id", Code: "ERR_INVALID_REQUEST",
		})
		return
	}

	marker, err := h.content.GetMarkerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status: "error", Message: "Marker not found", Code: "ERR_NOT_FOUND",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Marker retrieved", Data: marker,
	})
}

func (h *Handler) GetVideos(c *gin.Context) {
	videos, err := h.content.GetVideos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to load videos", Code: "ERR_DB_QUERY",
		})
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Videos retrieved", Data: videos,
	})
}

func (h *Handler) GetVideoDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ApiResponse{
			Status: "error", Message: "Invalid video id", Code: "ERR_INVALID_REQUEST",
		})
		return
	}

	video, err := h.content.GetVideoByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status: "error", Message: "Video not found", Code: "ERR_NOT_FOUND",
		})
		return
	}

	_ = h.content.IncrementViewCount(id)

	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Video retrieved", Data: video,
	})
}

func (h *Handler) GetFeatureSwitch(c *gin.Context) {
	feature := c.Param("feature")
	switchData, err := h.features.GetByName(feature)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ApiResponse{
			Status: "error", Message: "Feature not found", Code: "ERR_NOT_FOUND",
		})
		return
	}

	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Feature status retrieved", Data: switchData,
	})
}

func (h *Handler) GetAllFeatureSwitches(c *gin.Context) {
	features, err := h.features.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to load features", Code: "ERR_DB_QUERY",
		})
		return
	}
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Features retrieved", Data: features,
	})
}

func (h *Handler) EnableFeature(c *gin.Context) {
	feature := c.Param("feature")
	if err := h.features.SetStatus(feature, "active"); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to enable feature", Code: "ERR_DB_UPDATE",
		})
		return
	}
	_ = h.features.LogChange(feature, "activated", currentUserID(c))
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Feature enabled",
		Data: gin.H{"feature": feature, "status": "active"},
	})
}

func (h *Handler) DisableFeature(c *gin.Context) {
	feature := c.Param("feature")
	if err := h.features.SetStatus(feature, "inactive"); err != nil {
		c.JSON(http.StatusInternalServerError, models.ApiResponse{
			Status: "error", Message: "Failed to disable feature", Code: "ERR_DB_UPDATE",
		})
		return
	}
	_ = h.features.LogChange(feature, "deactivated", currentUserID(c))
	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Feature disabled",
		Data: gin.H{"feature": feature, "status": "inactive"},
	})
}

func (h *Handler) GetStats(c *gin.Context) {
	totalUsers, _ := h.users.Count()
	totalQuizzes, _ := h.quiz.CountQuestions()
	totalAttempts, _ := h.quiz.CountAttempts()

	c.JSON(http.StatusOK, models.ApiResponse{
		Status: "success", Message: "Statistics retrieved",
		Data: gin.H{
			"total_users":    totalUsers,
			"total_quizzes":  totalQuizzes,
			"total_attempts": totalAttempts,
		},
	})
}

func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return c.GetString("token")
}

func currentUserID(c *gin.Context) *int {
	if id, ok := c.Get("user_id"); ok {
		if v, ok := id.(int); ok && v > 0 {
			return &v
		}
	}
	return nil
}
