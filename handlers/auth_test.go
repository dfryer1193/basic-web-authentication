package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/dfryer1193/basic-web-authentication/models"
	"github.com/dfryer1193/basic-web-authentication/storage"
	"github.com/dfryer1193/basic-web-authentication/utils"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	mockStore := storage.NewInMemoryUserStore()
	password := "securepassword"
	passwordHash, _ := utils.HashPassword(password)
	mockStore.Set("testuser", models.User{
		Username:     "testuser",
		PasswordHash: passwordHash,
	})
	handler := UserAwareHandler{userStore: mockStore}

	tests := []struct {
		name         string
		method       string
		body         interface{}
		wantStatus   int
		wantResponse string
	}{
		{
			name:         "invalid method",
			method:       http.MethodGet,
			body:         nil,
			wantStatus:   http.StatusMethodNotAllowed,
			wantResponse: "Only POST method is allowed\n",
		},
		{
			name:       "invalid body",
			method:     http.MethodPost,
			body:       "invalid-json",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:         "non-existent user",
			method:       http.MethodPost,
			body:         models.Credentials{Username: "nonexistent", Password: "password"},
			wantStatus:   http.StatusUnauthorized,
			wantResponse: "Invalid username or password\n",
		},
		{
			name:         "incorrect password",
			method:       http.MethodPost,
			body:         models.Credentials{Username: "testuser", Password: "wrongpassword"},
			wantStatus:   http.StatusUnauthorized,
			wantResponse: "Invalid username or password\n",
		},
		{
			name:         "successful login",
			method:       http.MethodPost,
			body:         models.Credentials{Username: "testuser", Password: password},
			wantStatus:   http.StatusOK,
			wantResponse: "Login successful",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var reqBody []byte
			if tt.body != nil {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/login", bytes.NewReader(reqBody))
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}

			rec := httptest.NewRecorder()
			handler.LoginHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantStatus)
			}

			if tt.wantResponse != "" {
				respBody := rec.Body.String()
				if respBody != tt.wantResponse {
					t.Errorf("got response body %q, want %q", respBody, tt.wantResponse)
				}
			}
		})
	}
}

func TestRegisterHandler(t *testing.T) {
	mockStore := storage.NewInMemoryUserStore()
	handler := UserAwareHandler{userStore: mockStore}

	tests := []struct {
		name         string
		method       string
		body         interface{}
		wantStatus   int
		wantResponse string
		setup        func()
	}{
		{
			name:         "invalid method",
			method:       http.MethodGet,
			body:         nil,
			wantStatus:   http.StatusMethodNotAllowed,
			wantResponse: "Only POST method is allowed\n",
		},
		{
			name:         "invalid JSON payload",
			method:       http.MethodPost,
			body:         "invalid-json",
			wantStatus:   http.StatusBadRequest,
			wantResponse: "Invalid request body\n",
		},
		{
			name:         "valid registration",
			method:       http.MethodPost,
			body:         models.Credentials{Username: "newuser", Password: "password123"},
			wantStatus:   http.StatusCreated,
			wantResponse: "User registered successfully",
		},
		{
			name:   "duplicate user registration",
			method: http.MethodPost,
			body:   models.Credentials{Username: "existinguser", Password: "password123"},
			setup: func() {
				passwordHash, _ := utils.HashPassword("password123")
				mockStore.Set("existinguser", models.User{
					Username:     "existinguser",
					PasswordHash: passwordHash,
				})
			},
			wantStatus:   http.StatusBadRequest,
			wantResponse: "Username already exists\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup()
			}

			var reqBody []byte
			if tt.body != nil {
				reqBody, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/register", bytes.NewReader(reqBody))
			if tt.method == http.MethodPost {
				req.Header.Set("Content-Type", "application/json")
			}

			rec := httptest.NewRecorder()
			handler.RegisterHandler(rec, req)

			res := rec.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantStatus {
				t.Errorf("got status %v, want %v", res.StatusCode, tt.wantStatus)
			}

			respBody := rec.Body.String()
			if respBody != tt.wantResponse {
				t.Errorf("got response body %q, want %q", respBody, tt.wantResponse)
			}
		})
	}
}
