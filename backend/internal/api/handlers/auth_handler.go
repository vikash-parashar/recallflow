package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/recallflow/backend/internal/middleware"
	"github.com/recallflow/backend/internal/repositories"
	"github.com/recallflow/backend/internal/services"
)

type AuthHandler struct {
	authService *services.AuthService
	userRepo    *repositories.UserRepository
	orgRepo     *repositories.OrganizationRepository
}

func NewAuthHandler(authService *services.AuthService, userRepo *repositories.UserRepository, orgRepo *repositories.OrganizationRepository) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userRepo:    userRepo,
		orgRepo:     orgRepo,
	}
}

type RegisterRequest struct {
	OrganizationName string `json:"organization_name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
	Phone            string `json:"phone"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.OrganizationName == "" {
		RespondError(w, http.StatusBadRequest, "Email, password, and organization name are required")
		return
	}

	// Create organization and user
	token, user, err := h.authService.Register(
		req.OrganizationName,
		req.Email,
		req.Password,
		req.FirstName,
		req.LastName,
		req.Phone,
	)
	if err != nil {
		RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondJSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		RespondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	RespondJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserFromContext(r)
	if !ok {
		RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	user, err := h.userRepo.GetByID(claims.UserID)
	if err != nil {
		RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	RespondJSON(w, http.StatusOK, user)
}
