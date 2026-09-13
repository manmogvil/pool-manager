package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"lottery-pool-manager/models"
)

func TestContributionHandler_List_Success(t *testing.T) {
	mock := &mockContributionStore{
		getAllFn: func() ([]models.Contribution, error) {
			now := time.Now()
			return []models.Contribution{
				{ID: 1, ParticipantID: 1, GameID: 1, Month: 1, Year: 2024, Amount: 5.00, Paid: true, PaymentDate: &now, PaymentMethod: "CASH"},
			}, nil
		},
	}

	rr := doContributionRequest(t, mock, "GET", "/contributions", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}

	var result []models.Contribution
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 contribution, got %d", len(result))
	}
}

func TestContributionHandler_Create_Success(t *testing.T) {
	mock := &mockContributionStore{
		createFn: func(c *models.Contribution) (models.Contribution, error) {
			return models.Contribution{ID: 1, ParticipantID: c.ParticipantID, GameID: c.GameID, Amount: c.Amount}, nil
		},
	}

	rr := doContributionRequest(t, mock, "POST", "/contributions",
		`{"participant_id": 1, "game_id": 1, "month": 1, "year": 2024, "amount": 5.00, "paid": true, "payment_method": "CASH"}`)

	if rr.Code != http.StatusCreated {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusCreated)
	}
}

func TestContributionHandler_Create_MissingParticipantID(t *testing.T) {
	mock := &mockContributionStore{}
	rr := doContributionRequest(t, mock, "POST", "/contributions",
		`{"game_id": 1, "month": 1, "year": 2024, "amount": 5.00}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestContributionHandler_Create_MissingGameID(t *testing.T) {
	mock := &mockContributionStore{}
	rr := doContributionRequest(t, mock, "POST", "/contributions",
		`{"participant_id": 1, "month": 1, "year": 2024, "amount": 5.00}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestContributionHandler_Create_InvalidMonth(t *testing.T) {
	mock := &mockContributionStore{}
	rr := doContributionRequest(t, mock, "POST", "/contributions",
		`{"participant_id": 1, "game_id": 1, "month": 13, "year": 2024, "amount": 5.00}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestContributionHandler_Create_InvalidPaymentMethod(t *testing.T) {
	mock := &mockContributionStore{}
	rr := doContributionRequest(t, mock, "POST", "/contributions",
		`{"participant_id": 1, "game_id": 1, "month": 1, "year": 2024, "amount": 5.00, "payment_method": "CARD"}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestContributionHandler_Update_Success(t *testing.T) {
	mock := &mockContributionStore{
		updateFn: func(id int, c *models.Contribution) (models.Contribution, error) {
			return models.Contribution{ID: id, ParticipantID: c.ParticipantID, Amount: c.Amount}, nil
		},
	}

	rr := doContributionRequest(t, mock, "PUT", "/contributions/1",
		`{"participant_id": 1, "game_id": 1, "month": 1, "year": 2024, "amount": 10.00, "paid": true, "payment_method": "BIZUM"}`)

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestContributionHandler_Delete_Success(t *testing.T) {
	mock := &mockContributionStore{
		deleteFn: func(id int) error {
			return nil
		},
	}

	rr := doContributionRequest(t, mock, "DELETE", "/contributions/1", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestContributionHandler_Delete_NotFound(t *testing.T) {
	mock := &mockContributionStore{
		deleteFn: func(id int) error {
			return errors.New("contribution 99 not found")
		},
	}

	rr := doContributionRequest(t, mock, "DELETE", "/contributions/99", "")

	if rr.Code != http.StatusNotFound {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusNotFound)
	}
}
