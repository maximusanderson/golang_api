package handlers

import (
	"database/sql"
	"encoding/json"
	"golang_api/database"
	"golang_api/models"
	"log"
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	DB *database.DB
}

func NewAuthHandler(db *database.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	var user models.User
	var hashedPassword string
	err := h.DB.QueryRow(
		"SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Email, &hashedPassword, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	_, err = session.CreateNewSession(r, w, "public", user.Email, nil, nil)
	if err != nil {
		log.Printf("Error creating session: %v", err)
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.LoginResponse{
		Message: "Login successful",
		User:    &user,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	if sessionContainer == nil {
		http.Error(w, "No active session", http.StatusUnauthorized)
		return
	}

	if err := sessionContainer.RevokeSession(); err != nil {
		log.Printf("Error revoking session: %v", err)
		http.Error(w, "Error logging out", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successful",
	})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var user models.User
	err = h.DB.QueryRow(
		"INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id, email, created_at, updated_at",
		req.Email, string(hashedPassword),
	).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "User already exists or database error", http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.LoginResponse{
		Message: "Registration successful",
		User:    &user,
	})
}
