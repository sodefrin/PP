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

	SignupHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	if resp["Name"] != "testuser1" {
		t.Errorf("Expected name testuser1, got %v", resp["Name"])
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
	SignupHandler(w1, req1)
	if w1.Code != http.StatusCreated {
		t.Fatalf("Failed to create initial user: %d", w1.Code)
	}

	// Second creation (duplicate)
	req2 := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(body))
	w2 := httptest.NewRecorder()
	SignupHandler(w2, req2)

	if w2.Code != http.StatusConflict {
		t.Errorf("Expected status 409 for duplicate user, got %d", w2.Code)
	}
}
