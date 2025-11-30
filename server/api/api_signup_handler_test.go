package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignup(t *testing.T) {
	reqBody := map[string]string{
		"name":     "testuser1",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	if err := SignupHandler(testQueries)(w, req); err != nil {
		t.Fatalf("SignupHandler error: %v", err)
	}

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if resp["name"] != "testuser1" {
		t.Errorf("Expected name testuser1, got %v", resp["name"])
	}

	// Verify session cookie
	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == "session_id" {
			found = true
			if c.Value == "" {
				t.Error("Session cookie value is empty")
			}
			break
		}
	}
	if !found {
		t.Error("Session cookie not found")
	}
}

func TestSignupDuplicate(t *testing.T) {
	// Ensure user exists (from previous test or create new)
	// Since tests run in random order usually, let's create a specific user for this test
	reqBody := map[string]string{
		"name":     "dupuser",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)

	// First creation
	req1 := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	w1 := httptest.NewRecorder()
	if err := SignupHandler(testQueries)(w1, req1); err != nil {
		t.Fatalf("SignupHandler error: %v", err)
	}

	if w1.Code != http.StatusCreated {
		t.Fatalf("Failed to create initial user: %d", w1.Code)
	}

	// Second creation (duplicate)
	req2 := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	w2 := httptest.NewRecorder()
	if err := SignupHandler(testQueries)(w2, req2); err != nil {
		t.Fatalf("SignupHandler error: %v", err)
	}

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate user, got %d", w2.Code)
	}
}
