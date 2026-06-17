//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"
)

const defaultBaseURL = "http://localhost:8080"

func apiBase() string {
	if v := os.Getenv("API_BASE_URL"); v != "" {
		return v
	}
	return defaultBaseURL
}

func TestIntegrationHealth(t *testing.T) {
	resp, err := http.Get(apiBase() + "/health")
	if err != nil {
		t.Skipf("API not running: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health status = %d", resp.StatusCode)
	}
}

func TestIntegrationAuthAndQuizFlow(t *testing.T) {
	base := apiBase()

	signInBody, _ := json.Marshal(map[string]string{
		"email":    "editor@sasivision.com",
		"password": "Sasivision123",
	})
	signInResp, err := http.Post(base+"/api/auth/sign-in", "application/json", bytes.NewReader(signInBody))
	if err != nil {
		t.Skipf("API not running: %v", err)
	}
	defer signInResp.Body.Close()
	if signInResp.StatusCode != http.StatusOK {
		t.Fatalf("sign-in status = %d", signInResp.StatusCode)
	}

	var signIn struct {
		Data struct {
			Token string `json:"token"`
			Role  string `json:"role"`
		} `json:"data"`
	}
	if err := json.NewDecoder(signInResp.Body).Decode(&signIn); err != nil {
		t.Fatal(err)
	}
	if signIn.Data.Role != "editor" {
		t.Fatalf("role = %q, want editor", signIn.Data.Role)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest(http.MethodGet, base+"/api/admin/quiz/categories", nil)
	req.Header.Set("Authorization", "Bearer "+signIn.Data.Token)
	catResp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer catResp.Body.Close()
	if catResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(catResp.Body)
		t.Fatalf("categories status = %d body=%s", catResp.StatusCode, body)
	}

	pubResp, err := http.Get(base + "/api/content/videos")
	if err != nil {
		t.Fatal(err)
	}
	defer pubResp.Body.Close()
	if pubResp.StatusCode != http.StatusOK {
		t.Fatalf("videos status = %d", pubResp.StatusCode)
	}
}

func TestIntegrationQuizSubmit(t *testing.T) {
	base := apiBase()

	signInBody, _ := json.Marshal(map[string]string{
		"email":    "demo@sasivision.com",
		"password": "Sasivision123",
	})
	signInResp, err := http.Post(base+"/api/auth/sign-in", "application/json", bytes.NewReader(signInBody))
	if err != nil {
		t.Skipf("API not running: %v", err)
	}
	defer signInResp.Body.Close()
	if signInResp.StatusCode != http.StatusOK {
		t.Fatalf("sign-in status = %d", signInResp.StatusCode)
	}

	var signIn struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(signInResp.Body).Decode(&signIn); err != nil {
		t.Fatal(err)
	}

	qResp, err := http.Get(base + "/api/quiz/questions/Post-Test")
	if err != nil {
		t.Fatal(err)
	}
	defer qResp.Body.Close()
	if qResp.StatusCode != http.StatusOK {
		t.Fatalf("questions status = %d", qResp.StatusCode)
	}

	var questions struct {
		Data []struct {
			ID      int `json:"id"`
			Answers []struct {
				AnswerKey  string `json:"answer_key"`
				IsCorrect  bool   `json:"is_correct"`
			} `json:"answers"`
			Type string `json:"type"`
		} `json:"data"`
	}
	if err := json.NewDecoder(qResp.Body).Decode(&questions); err != nil {
		t.Fatal(err)
	}
	if len(questions.Data) < 5 {
		t.Fatalf("expected >= 5 questions, got %d", len(questions.Data))
	}

	answers := make([]map[string]interface{}, 0, len(questions.Data))
	correct := 0
	for _, q := range questions.Data {
		ans := "essay answer"
		if q.Type == "multiple_choice" {
			for _, opt := range q.Answers {
				if opt.IsCorrect {
					ans = opt.AnswerKey
					correct++
					break
				}
			}
		}
		answers = append(answers, map[string]interface{}{
			"quiz_id":  q.ID,
			"category": "Post-Test",
			"type":     q.Type,
			"answers":  ans,
		})
	}

	now := time.Now().UTC().Format(time.RFC3339)
	submitBody, _ := json.Marshal(map[string]interface{}{
		"email":       "demo@sasivision.com",
		"category_id": 1,
		"correct":     correct,
		"total":       len(questions.Data),
		"score":       80,
		"start_time":  now,
		"end_time":    now,
		"finish_date": time.Now().Format("2006-01-02"),
		"answers":     answers,
	})

	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest(http.MethodPost, base+"/api/quiz/attempts", bytes.NewReader(submitBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+signIn.Data.Token)
	submitResp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer submitResp.Body.Close()
	if submitResp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(submitResp.Body)
		t.Fatalf("submit status = %d body=%s", submitResp.StatusCode, body)
	}
}
