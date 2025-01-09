package handlers

import (
	"encoding/json"
	"github.com/dfryer1193/basic-web-authentication/models"
	"github.com/dfryer1193/basic-web-authentication/storage"
	"github.com/dfryer1193/basic-web-authentication/utils"
	"net/http"
)

type UserAwareHandler struct {
	cookieName string
	userStore  storage.UserStore
}

// NewUserAwareHandler initializes and returns a UserAwareHandler with the provided UserStore dependency.
func NewUserAwareHandler(cookieName string, userStore storage.UserStore) UserAwareHandler {
	return UserAwareHandler{
		cookieName: cookieName,
		userStore:  userStore,
	}
}

// LoginHandler processes user login requests, validates credentials, and sets a session cookie for successful authentication.
func (handler UserAwareHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, ok := handler.userStore.Get(credentials.Username)
	if !ok || !utils.CheckPasswordHash(credentials.Password, user.PasswordHash) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Authentication successful
	http.SetCookie(w, &http.Cookie{Name: "session", Value: credentials.Username, Path: "/"})
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}

// RegisterHandler handles user registration by accepting credentials via a POST request, hashing the password, and storing the user.
func (handler UserAwareHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var credentials models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, exists := handler.userStore.Get(credentials.Username)
	if exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	hash, err := utils.HashPassword(credentials.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	handler.userStore.Set(credentials.Username, models.User{
		Username:     credentials.Username,
		PasswordHash: hash,
	})

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User registered successfully"))
}

// WelcomeHandler validates the user's session cookie and responds with a welcome message if the user is authenticated.
func (handler UserAwareHandler) WelcomeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(handler.cookieName)
	if err != nil || cookie.Value == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome, " + cookie.Value))
}
