package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignin(t *testing.T) {
	// Create user first
	signupBody := map[string]string{
		"name":     "signinuser",
		"password": "password123",
	}
	sBody, _ := json.Marshal(signupBody)
	sReq := httptest.NewRequest(http.MethodPost, "/api/signup", bytes.NewBuffer(sBody))
	sW := httptest.NewRecorder()
	if err := SignupHandler(testQueries)(sW, sReq); err != nil {
		t.Fatalf("SignupHandler error: %v", err)
	}

	// Test Signin
	signinBody := map[string]string{
		"name":     "signinuser",
		"password": "password123",
	}
	body, _ := json.Marshal(signinBody)
	req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	if err := SigninHandler(testQueries)(w, req); err != nil {
		t.Fatalf("SigninHandler error: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
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

func TestSigninInvalid(t *testing.T) {
	signinBody := map[string]string{
		"name":     "nonexistent",
		"password": "password123",
	}
	body, _ := json.Marshal(signinBody)
	req := httptest.NewRequest(http.MethodPost, "/api/signin", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	if err := SigninHandler(testQueries)(w, req); err != nil {
		t.Fatalf("SigninHandler error: %v", err)
	}

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}
