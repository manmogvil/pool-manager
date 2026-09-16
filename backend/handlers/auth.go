package handlers

import (
	"encoding/json"
	"net/http"

	"lottery-pool-manager/models"
	"lottery-pool-manager/services"
	"lottery-pool-manager/utils"
)

type AuthStore interface {
	CreateUser(name, email, passwordHash, role string) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(id int) (models.User, error)
	ActivateUser(id int) error
	DeactivateUser(id int) error
	UpdateUser(id int, name, email, role string) (models.User, error)
	GetAllUsers() ([]models.User, error)
}

type AuthHandler struct {
	store AuthStore
	auth  *services.AuthService
}

func NewAuthHandler(s AuthStore) *AuthHandler {
	return &AuthHandler{
		store: s,
		auth:  services.NewAuthService(),
	}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		utils.RespondError(w, http.StatusBadRequest, "Name, email and password are required")
		return
	}

	if req.Role == "" {
		req.Role = "user"
	}

	if req.Role != "admin" && req.Role != "user" {
		utils.RespondError(w, http.StatusBadRequest, "Role must be 'admin' or 'user'")
		return
	}

	existing, _ := h.store.GetUserByEmail(req.Email)
	if existing.ID != 0 {
		utils.RespondError(w, http.StatusConflict, "Email already registered")
		return
	}

	hash, err := h.auth.HashPassword(req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error hashing password")
		return
	}

	user, err := h.store.CreateUser(req.Name, req.Email, hash, req.Role)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error creating user")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "Registration successful. Waiting for admin approval.",
		"user":    user,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.RespondError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	user, err := h.store.GetUserByEmail(req.Email)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	if !user.Active {
		utils.RespondError(w, http.StatusForbidden, "Account not activated. Wait for admin approval.")
		return
	}

	if !h.auth.CheckPassword(req.Password, user.PasswordHash) {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := h.auth.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error generating token")
		return
	}

	utils.RespondJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	user, err := h.store.GetUserByID(userID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "User not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetAllUsers()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching users")
		return
	}

	utils.RespondJSON(w, http.StatusOK, users)
}

func (h *AuthHandler) ListUserNames(w http.ResponseWriter, r *http.Request) {
	users, err := h.store.GetAllUsers()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching users")
		return
	}

	type userName struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	var result []userName
	for _, u := range users {
		if u.Active {
			result = append(result, userName{ID: u.ID, Name: u.Name})
		}
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

func (h *AuthHandler) ActivateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "User ID required")
		return
	}

	id, err := utils.ParseID(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.store.ActivateUser(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "User activated"})
}

func (h *AuthHandler) DeactivateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "User ID required")
		return
	}

	id, err := utils.ParseID(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.store.DeactivateUser(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "User deactivated"})
}

type updateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "User ID required")
		return
	}

	id, err := utils.ParseID(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" || req.Email == "" || req.Role == "" {
		utils.RespondError(w, http.StatusBadRequest, "Name, email and role are required")
		return
	}

	if req.Role != "admin" && req.Role != "user" {
		utils.RespondError(w, http.StatusBadRequest, "Role must be 'admin' or 'user'")
		return
	}

	user, err := h.store.UpdateUser(id, req.Name, req.Email, req.Role)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, user)
}
