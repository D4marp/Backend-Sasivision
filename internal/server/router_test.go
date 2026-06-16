package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/sasivision/backend/internal/config"
	"github.com/sasivision/backend/internal/handlers"
	"github.com/sasivision/backend/internal/utils"
)

func testRouter(t *testing.T) (*gin.Engine, sqlmock.Sqlmock, *sql.DB) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("sqlmock.New: %v", err)
	}
	cfg := &config.Config{
		JWTSecret:    "api-test-secret",
		JWTExpiry:    time.Hour,
		CORSOrigin:   "*",
		RateLimitRPS: 1000,
	}
	h := handlers.New(db, cfg)
	return SetupRouter(h, cfg), mock, db
}

func bearerToken(t *testing.T, cfg *config.Config, userID int, email, role string) string {
	t.Helper()
	token, err := utils.GenerateToken(userID, email, role, cfg.JWTSecret, cfg.JWTExpiry)
	if err != nil {
		t.Fatal(err)
	}
	return token
}

func TestHealthEndpoint(t *testing.T) {
	router, _, db := testRouter(t)
	defer db.Close()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
}

func TestGetQuizCategoriesPublic(t *testing.T) {
	router, mock, db := testRouter(t)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "name", "slug", "description", "display_order", "is_active", "created_at",
	}).AddRow(1, "Motif Dasar", "motif-dasar", "Deskripsi", 1, true, now)
	mock.ExpectQuery("SELECT id, name, slug, description, display_order, is_active, created_at").
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/quiz/categories", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestAdminQuizQuestionsRequiresEditorRole(t *testing.T) {
	router, _, db := testRouter(t)
	defer db.Close()

	cfg := &config.Config{JWTSecret: "api-test-secret", JWTExpiry: time.Hour, RateLimitRPS: 1000}
	userToken := bearerToken(t, cfg, 2, "mahasiswa@test.com", "user")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/quiz/questions?category_id=1", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestAdminQuizQuestionsEditorSuccess(t *testing.T) {
	router, mock, db := testRouter(t)
	defer db.Close()

	cfg := &config.Config{JWTSecret: "api-test-secret", JWTExpiry: time.Hour, RateLimitRPS: 1000}
	editorToken := bearerToken(t, cfg, 3, "editor@sasivision.com", "editor")

	now := time.Now()
	qRows := sqlmock.NewRows([]string{
		"id", "category_id", "type", "question", "image_url", "sequence_order", "created_at",
	}).AddRow(10, 1, "multiple_choice", "Apa warna motif?", nil, 1, now)
	mock.ExpectQuery("SELECT id, category_id, type, question, image_url, sequence_order, created_at").
		WithArgs(1).
		WillReturnRows(qRows)

	aRows := sqlmock.NewRows([]string{
		"id", "quiz_id", "answer_key", "answer_text", "is_correct", "created_at",
	}).AddRow(1, 10, "A", "Merah", true, now)
	mock.ExpectQuery("SELECT id, quiz_id, answer_key, answer_text, is_correct, created_at").
		WithArgs(10).
		WillReturnRows(aRows)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/quiz/questions?category_id=1", nil)
	req.Header.Set("Authorization", "Bearer "+editorToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestGetAttemptDetailsForbiddenForOtherUser(t *testing.T) {
	router, mock, db := testRouter(t)
	defer db.Close()

	cfg := &config.Config{JWTSecret: "api-test-secret", JWTExpiry: time.Hour, RateLimitRPS: 1000}
	userToken := bearerToken(t, cfg, 5, "user-a@test.com", "user")

	ownerRows := sqlmock.NewRows([]string{"email"}).AddRow("user-b@test.com")
	mock.ExpectQuery("SELECT u.email FROM quiz_attempts qa").
		WithArgs(99).
		WillReturnRows(ownerRows)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/quiz/attempts/99/details", nil)
	req.Header.Set("Authorization", "Bearer "+userToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestCreateQuizCategoryEditor(t *testing.T) {
	router, mock, db := testRouter(t)
	defer db.Close()

	cfg := &config.Config{JWTSecret: "api-test-secret", JWTExpiry: time.Hour, RateLimitRPS: 1000}
	editorToken := bearerToken(t, cfg, 3, "editor@sasivision.com", "editor")

	mock.ExpectExec("INSERT INTO quiz_categories").
		WithArgs("Kategori Baru", "kategori-baru", "Deskripsi", 2, true).
		WillReturnResult(sqlmock.NewResult(12, 1))

	body, _ := json.Marshal(map[string]interface{}{
		"name":          "Kategori Baru",
		"description":   "Deskripsi",
		"display_order": 2,
		"is_active":     true,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/quiz/categories", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+editorToken)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateVideoEditor(t *testing.T) {
	router, mock, db := testRouter(t)
	defer db.Close()

	cfg := &config.Config{JWTSecret: "api-test-secret", JWTExpiry: time.Hour, RateLimitRPS: 1000}
	editorToken := bearerToken(t, cfg, 3, "editor@sasivision.com", "editor")

	mock.ExpectExec("INSERT INTO videos").
		WillReturnResult(sqlmock.NewResult(5, 1))

	body, _ := json.Marshal(map[string]interface{}{
		"title":       "Video Edukasi",
		"description": "Materi pembelajaran",
		"source":      "local",
		"video_url":   "videos/sample.mp4",
		"thumbnail":   "videos/thumbnails/sample.jpg",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/videos", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+editorToken)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminUsersBlockedForEditor(t *testing.T) {
	router, _, db := testRouter(t)
	defer db.Close()

	cfg := &config.Config{JWTSecret: "api-test-secret", JWTExpiry: time.Hour, RateLimitRPS: 1000}
	editorToken := bearerToken(t, cfg, 3, "editor@sasivision.com", "editor")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+editorToken)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", w.Code)
	}
}

func TestGetMarkersPublic(t *testing.T) {
	router, mock, db := testRouter(t)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id", "title", "slug", "description", "image_file", "audio_file", "model_path", "sentences", "display_order", "created_at",
	}).AddRow(1, "Naga Balimbur", "naga-balimbur", "Desc", "markers/naga.png", "", "", "[]", 1, now)
	mock.ExpectQuery("SELECT id, title, slug, description, image_file, audio_file, model_path, sentences, display_order, created_at FROM markers").
		WillReturnRows(rows)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/content/markers", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}

func TestCreateMarkerEditor(t *testing.T) {
	router, mock, db := testRouter(t)
	defer db.Close()

	cfg := &config.Config{JWTSecret: "api-test-secret", JWTExpiry: time.Hour, RateLimitRPS: 1000}
	editorToken := bearerToken(t, cfg, 3, "editor@sasivision.com", "editor")

	mock.ExpectExec("INSERT INTO markers").
		WillReturnResult(sqlmock.NewResult(8, 1))

	body, _ := json.Marshal(map[string]interface{}{
		"title":       "Motif Baru",
		"description": "Materi AR",
		"image_file":  "markers/baru.png",
		"audio_file":  "audio/baru.mp3",
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/admin/markers", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+editorToken)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d body=%s", w.Code, w.Body.String())
	}
}
