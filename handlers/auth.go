package handlers

import (
	"encoding/json"
	"net/http"
)

var userStore = storage.NewInMemoryUserStore()

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, ok := userStore.Get(credentials.Username)
	if !ok || !utils.CheckPasswordHash(credentials.Password, user.PasswordHash) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Authentication successful
	http.SetCookie(w, &http.Cookie{Name: "session", Value: credentials.Username, Path: "/"})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// RegisterHandler handles user registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hash, err := utils.HashPassword(credentials.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	userStore.Set(credentials.Username, models.User{
		Username:     credentials.Username,
		PasswordHash: hash,
	})

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

// WelcomeHandler handles welcome page
func WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome, " + cookie.Value))
}
