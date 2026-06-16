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
