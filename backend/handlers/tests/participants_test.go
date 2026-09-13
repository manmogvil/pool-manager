package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"lottery-pool-manager/models"
)

func TestParticipantHandler_List_Success(t *testing.T) {
	mock := &mockParticipantStore{
		getAllFn: func() ([]models.Participant, error) {
			return []models.Participant{
				{ID: 1, Name: "Alice", Email: "alice@example.com", Active: true},
			}, nil
		},
	}

	rr := doParticipantRequest(t, mock, "GET", "/participants", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}

	var result []models.Participant
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 participant, got %d", len(result))
	}
}

func TestParticipantHandler_List_Error(t *testing.T) {
	mock := &mockParticipantStore{
		getAllFn: func() ([]models.Participant, error) {
			return nil, errors.New("database error")
		},
	}

	rr := doParticipantRequest(t, mock, "GET", "/participants", "")

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusInternalServerError)
	}
}

func TestParticipantHandler_GetByID_Success(t *testing.T) {
	mock := &mockParticipantStore{
		getByIDFn: func(id int) (models.Participant, error) {
			return models.Participant{ID: id, Name: "Alice", Email: "alice@example.com", Active: true}, nil
		},
	}

	rr := doParticipantRequest(t, mock, "GET", "/participants/1", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestParticipantHandler_GetByID_InvalidID(t *testing.T) {
	mock := &mockParticipantStore{}
	rr := doParticipantRequest(t, mock, "GET", "/participants/abc", "")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestParticipantHandler_GetByID_NotFound(t *testing.T) {
	mock := &mockParticipantStore{
		getByIDFn: func(id int) (models.Participant, error) {
			return models.Participant{}, errors.New("participant 99 not found")
		},
	}

	rr := doParticipantRequest(t, mock, "GET", "/participants/99", "")

	if rr.Code != http.StatusNotFound {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusNotFound)
	}
}

func TestParticipantHandler_Create_Success(t *testing.T) {
	mock := &mockParticipantStore{
		createFn: func(name, email string) (models.Participant, error) {
			return models.Participant{ID: 1, Name: name, Email: email, Active: true}, nil
		},
	}

	rr := doParticipantRequest(t, mock, "POST", "/participants",
		`{"name": "Alice", "email": "alice@example.com"}`)

	if rr.Code != http.StatusCreated {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusCreated)
	}
}

func TestParticipantHandler_Create_MissingName(t *testing.T) {
	mock := &mockParticipantStore{}
	rr := doParticipantRequest(t, mock, "POST", "/participants",
		`{"email": "alice@example.com"}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestParticipantHandler_Create_MissingEmail(t *testing.T) {
	mock := &mockParticipantStore{}
	rr := doParticipantRequest(t, mock, "POST", "/participants",
		`{"name": "Alice"}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestParticipantHandler_Create_InvalidJSON(t *testing.T) {
	mock := &mockParticipantStore{}
	rr := doParticipantRequest(t, mock, "POST", "/participants", "invalid")

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestParticipantHandler_Update_Success(t *testing.T) {
	mock := &mockParticipantStore{
		updateFn: func(id int, name, email string) (models.Participant, error) {
			return models.Participant{ID: id, Name: name, Email: email, Active: true}, nil
		},
	}

	rr := doParticipantRequest(t, mock, "PUT", "/participants/1",
		`{"name": "Updated Name", "email": "updated@example.com"}`)

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestParticipantHandler_Deactivate_Success(t *testing.T) {
	mock := &mockParticipantStore{
		deactivateFn: func(id int) (models.Participant, error) {
			return models.Participant{ID: id, Name: "Alice", Active: false}, nil
		},
	}

	rr := doParticipantRequest(t, mock, "PUT", "/participants/1/deactivate", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestParticipantHandler_Activate_Success(t *testing.T) {
	mock := &mockParticipantStore{
		activateFn: func(id int) (models.Participant, error) {
			return models.Participant{ID: id, Name: "Alice", Active: true}, nil
		},
	}

	rr := doParticipantRequest(t, mock, "PUT", "/participants/1/activate", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}
