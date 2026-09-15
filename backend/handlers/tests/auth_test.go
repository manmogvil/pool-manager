package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"lottery-pool-manager/handlers"
	"lottery-pool-manager/models"
	"lottery-pool-manager/services"
)

func setupAuthRouter(h *handlers.AuthHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
	r.Get("/auth/me", h.GetMe)
	r.Get("/auth/users", h.ListUsers)
	r.Put("/auth/users/{id}", h.UpdateUser)
	r.Put("/auth/users/{id}/activate", h.ActivateUser)
	r.Put("/auth/users/{id}/deactivate", h.DeactivateUser)
	return r
}

// injectUserID is middleware that sets user_id in context for GetMe tests
func injectUserID(userID int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func doAuthRequest(t *testing.T, handler http.Handler, method, url, body string) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, url, bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// --- Register Tests ---

func TestRegister_Success(t *testing.T) {
	mock := &mockAuthStore{
		getByEmailFn: func(email string) (models.User, error) {
			return models.User{}, errors.New("not found")
		},
		createUserFn: func(name, email, passwordHash, role string) (models.User, error) {
			return models.User{ID: 1, Name: name, Email: email, Role: role, Active: false}, nil
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "POST", "/auth/register",
		`{"name":"Test","email":"test@example.com","password":"pass123"}`)

	if rr.Code != http.StatusCreated {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusCreated)
	}
}

func TestRegister_MissingFields(t *testing.T) {
	mock := &mockAuthStore{}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "POST", "/auth/register",
		`{"name":"Test"}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	mock := &mockAuthStore{
		getByEmailFn: func(email string) (models.User, error) {
			return models.User{ID: 1, Email: email}, nil
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "POST", "/auth/register",
		`{"name":"Test","email":"existing@example.com","password":"pass123"}`)

	if rr.Code != http.StatusConflict {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusConflict)
	}
}

func TestRegister_InvalidRole(t *testing.T) {
	mock := &mockAuthStore{
		getByEmailFn: func(email string) (models.User, error) {
			return models.User{}, errors.New("not found")
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "POST", "/auth/register",
		`{"name":"Test","email":"test@example.com","password":"pass123","role":"superadmin"}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

// --- Login Tests ---

func TestLogin_Success(t *testing.T) {
	auth := services.NewAuthService()
	hash, _ := auth.HashPassword("password123")

	mock := &mockAuthStore{
		getByEmailFn: func(email string) (models.User, error) {
			return models.User{ID: 1, Name: "Admin", Email: email, PasswordHash: hash, Role: "admin", Active: true}, nil
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "POST", "/auth/login",
		`{"email":"admin@example.com","password":"password123"}`)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusOK)
	}

	var result handlers.AuthResponse
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if result.Token == "" {
		t.Error("expected non-empty token")
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mock := &mockAuthStore{
		getByEmailFn: func(email string) (models.User, error) {
			return models.User{}, errors.New("not found")
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "POST", "/auth/login",
		`{"email":"wrong@example.com","password":"wrong"}`)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusUnauthorized)
	}
}

func TestLogin_InactiveUser(t *testing.T) {
	auth := services.NewAuthService()
	hash, _ := auth.HashPassword("password123")

	mock := &mockAuthStore{
		getByEmailFn: func(email string) (models.User, error) {
			return models.User{ID: 1, Email: email, PasswordHash: hash, Active: false}, nil
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "POST", "/auth/login",
		`{"email":"inactive@example.com","password":"password123"}`)

	if rr.Code != http.StatusForbidden {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusForbidden)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	auth := services.NewAuthService()
	hash, _ := auth.HashPassword("correctpassword")

	mock := &mockAuthStore{
		getByEmailFn: func(email string) (models.User, error) {
			return models.User{ID: 1, Email: email, PasswordHash: hash, Active: true}, nil
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "POST", "/auth/login",
		`{"email":"user@example.com","password":"wrongpassword"}`)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusUnauthorized)
	}
}

// --- ListUsers Tests ---

func TestListUsers_Success(t *testing.T) {
	mock := &mockAuthStore{
		getAllUsersFn: func() ([]models.User, error) {
			return []models.User{
				{ID: 1, Name: "Admin", Role: "admin", Active: true},
				{ID: 2, Name: "Alice", Role: "user", Active: true},
			}, nil
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "GET", "/auth/users", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusOK)
	}

	var result []models.User
	json.NewDecoder(rr.Body).Decode(&result)
	if len(result) != 2 {
		t.Errorf("expected 2 users, got %d", len(result))
	}
}

func TestListUsers_Error(t *testing.T) {
	mock := &mockAuthStore{
		getAllUsersFn: func() ([]models.User, error) {
			return nil, errors.New("db error")
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "GET", "/auth/users", "")

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusInternalServerError)
	}
}

// --- UpdateUser Tests ---

func TestUpdateUser_Success(t *testing.T) {
	mock := &mockAuthStore{
		updateUserFn: func(id int, name, email, role string) (models.User, error) {
			return models.User{ID: id, Name: name, Email: email, Role: role, Active: true}, nil
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "PUT", "/auth/users/1",
		`{"name":"Updated","email":"updated@example.com","role":"admin"}`)

	if rr.Code != http.StatusOK {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestUpdateUser_MissingFields(t *testing.T) {
	mock := &mockAuthStore{}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "PUT", "/auth/users/1",
		`{"name":"Updated"}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestUpdateUser_InvalidRole(t *testing.T) {
	mock := &mockAuthStore{}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "PUT", "/auth/users/1",
		`{"name":"Updated","email":"e@e.com","role":"superadmin"}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	mock := &mockAuthStore{
		updateUserFn: func(id int, name, email, role string) (models.User, error) {
			return models.User{}, errors.New("user 99 not found")
		},
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "PUT", "/auth/users/99",
		`{"name":"Updated","email":"e@e.com","role":"user"}`)

	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusNotFound)
	}
}

// --- ActivateUser Tests ---

func TestActivateUser_Success(t *testing.T) {
	mock := &mockAuthStore{
		activateFn: func(id int) error { return nil },
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "PUT", "/auth/users/1/activate", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestActivateUser_NotFound(t *testing.T) {
	mock := &mockAuthStore{
		activateFn: func(id int) error { return errors.New("user 99 not found") },
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "PUT", "/auth/users/99/activate", "")

	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusNotFound)
	}
}

// --- DeactivateUser Tests ---

func TestDeactivateUser_Success(t *testing.T) {
	mock := &mockAuthStore{
		deactivateFn: func(id int) error { return nil },
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "PUT", "/auth/users/1/deactivate", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestDeactivateUser_NotFound(t *testing.T) {
	mock := &mockAuthStore{
		deactivateFn: func(id int) error { return errors.New("user 99 not found") },
	}
	h := handlers.NewAuthHandler(mock)
	rr := doAuthRequest(t, setupAuthRouter(h), "PUT", "/auth/users/99/deactivate", "")

	if rr.Code != http.StatusNotFound {
		t.Errorf("status = %v, want %v", rr.Code, http.StatusNotFound)
	}
}
