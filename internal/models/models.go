package models

import "time"

// User represents a user account
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FullName  string    `json:"full_name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserSession represents an active session for auto-login
type UserSession struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// QuizCategory represents a quiz topic
type QuizCategory struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Description  string    `json:"description"`
	DisplayOrder int       `json:"display_order"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// Quiz represents a quiz question
type Quiz struct {
	ID            int        `json:"id"`
	CategoryID    int        `json:"category_id"`
	Type          string     `json:"type"` // multiple_choice or essay
	Question      string     `json:"question"`
	ImageURL      *string    `json:"image_url"`
	SequenceOrder int        `json:"sequence_order"`
	Answers       []QuizAnswer `json:"options"`
	CreatedAt     time.Time  `json:"created_at"`
}

// QuizAnswer represents an answer option
type QuizAnswer struct {
	ID        int       `json:"id"`
	QuizID    int       `json:"quiz_id"`
	AnswerKey *string   `json:"answer_key"` // A-E
	AnswerText string   `json:"answer_text"`
	IsCorrect bool      `json:"is_correct"`
	CreatedAt time.Time `json:"created_at"`
}

// QuizAttempt represents a user's quiz submission
type QuizAttempt struct {
	ID          int                `json:"id"`
	UserID      int                `json:"user_id"`
	Email       string             `json:"email"`
	CategoryID  int                `json:"category_id"`
	CategoryName string            `json:"category_name"`
	CorrectCount int               `json:"correct"`
	TotalCount  int                `json:"total"`
	Score       int                `json:"score"`
	StartTime   time.Time          `json:"start_time"`
	EndTime     time.Time          `json:"end_time"`
	FinishDate  string             `json:"finish_date"`
	Details     []AttemptDetail    `json:"answers"`
	CreatedAt   time.Time          `json:"created_at"`
}

// AttemptDetail represents a user's answer to a quiz question
type AttemptDetail struct {
	ID              int       `json:"id"`
	QuizAttemptID   int       `json:"quiz_attempt_id"`
	QuizID          int       `json:"quiz_id"`
	Type            string    `json:"type"` // multiple_choice or essay
	UserAnswer      string    `json:"answers"`
	CreatedAt       time.Time `json:"created_at"`
}

// Marker represents a Sasirangan motif
type Marker struct {
	ID          int       `json:"id"`
	Title       string    `json:"markerTitle"`
	Slug        string    `json:"slug"`
	Description string    `json:"engDescription"`
	ImageFile   string    `json:"markerFile"`
	AudioFile   string    `json:"engAudioFile"`
	ModelPath   *string   `json:"modelPath"`
	Sentences   *string   `json:"sentences"`
	DisplayOrder int      `json:"display_order"`
	CreatedAt   time.Time `json:"created_at"`
}

// Video represents an educational video
type Video struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Slug             string    `json:"slug"`
	Description      string    `json:"desc"`
	Source           string    `json:"source"`
	VideoURL         string    `json:"video_url"`
	Thumbnail        string    `json:"thumbnail"`
	DiscussionFormURL *string  `json:"discussion_form_url"`
	ViewCount        int       `json:"view_count"`
	DisplayOrder     int       `json:"display_order"`
	CreatedAt        time.Time `json:"created_at"`
}

// FeatureSwitch represents a runtime feature toggle
type FeatureSwitch struct {
	ID          int       `json:"id"`
	FeatureName string    `json:"feature_name"`
	Status      string    `json:"status"` // active or inactive
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SignInRequest represents a sign-in request
type SignInRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// SignUpRequest represents a sign-up request
type SignUpRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name"`
}

// QuizSubmitRequest represents quiz submission payload
type QuizSubmitRequest struct {
	Email      string              `json:"email" binding:"required,email"`
	CategoryID int                 `json:"category_id" binding:"required"`
	Correct    int                 `json:"correct"`
	Total      int                 `json:"total"`
	Score      int                 `json:"score"`
	StartTime  string              `json:"start_time"`
	EndTime    string              `json:"end_time"`
	FinishDate string              `json:"finish_date"`
	Answers    []QuizAnswerPayload `json:"answers"`
}

// QuizAnswerPayload represents a single quiz answer submission
type QuizAnswerPayload struct {
	QuizID   int    `json:"quiz_id"`
	Category string `json:"category"`
	Type     string `json:"type"`
	Answers  string `json:"answers"`
}

// SignInResponse represents a sign-in response
type SignInResponse struct {
	Status string `json:"status"`
	Email  string `json:"email"`
	Token  string `json:"token"`
}

// ApiResponse represents a standardized API response
type ApiResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Code    string      `json:"code,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// --- Admin CRUD request payloads ---

// VideoRequest is the create/update payload for a video
type VideoRequest struct {
	Title             string  `json:"title" binding:"required"`
	Slug              string  `json:"slug"`
	Description       string  `json:"description"`
	Source            string  `json:"source"`
	VideoURL          string  `json:"video_url" binding:"required"`
	Thumbnail         string  `json:"thumbnail" binding:"required"`
	DiscussionFormURL *string `json:"discussion_form_url"`
	DisplayOrder      int     `json:"display_order"`
}

// MarkerRequest is the create/update payload for a marker
type MarkerRequest struct {
	Title        string   `json:"title" binding:"required"`
	Slug         string   `json:"slug"`
	Description  string   `json:"description"`
	ImageFile    string   `json:"image_file" binding:"required"`
	AudioFile    string   `json:"audio_file" binding:"required"`
	ModelPath    *string  `json:"model_path"`
	Sentences    []string `json:"sentences"`
	DisplayOrder int      `json:"display_order"`
}

// QuizCategoryRequest is the create/update payload for a quiz category
type QuizCategoryRequest struct {
	Name         string `json:"name" binding:"required"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
	IsActive     bool   `json:"is_active"`
}

// QuizAnswerRequest is a single answer option for a question
type QuizAnswerRequest struct {
	AnswerKey  string `json:"answer_key"`
	AnswerText string `json:"answer_text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
}

// QuizQuestionRequest is the create/update payload for a quiz question + answers
type QuizQuestionRequest struct {
	CategoryID    int                 `json:"category_id" binding:"required"`
	Type          string              `json:"type" binding:"required"`
	Question      string              `json:"question" binding:"required"`
	ImageURL      *string             `json:"image_url"`
	SequenceOrder int                 `json:"sequence_order"`
	Answers       []QuizAnswerRequest `json:"answers"`
}

// UpdateRoleRequest changes a user's role
type UpdateRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=user editor admin"`
}

// AnalyticsEventRequest records a client-side analytics event
type AnalyticsEventRequest struct {
	EventType  string                 `json:"event_type" binding:"required"`
	EntityType string                 `json:"entity_type"`
	EntityID   int                    `json:"entity_id"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// CategoryCount is a generic label/count pair for charts
type CategoryCount struct {
	Label string `json:"label"`
	Count int    `json:"count"`
}

// TimePoint is a date/value pair for time-series charts
type TimePoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}
