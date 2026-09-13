package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"lottery-pool-manager/models"
)

func TestTicketHandler_List_Success(t *testing.T) {
	mock := &mockTicketStore{
		getAllFn: func() ([]models.Ticket, error) {
			return []models.Ticket{
				{ID: 1, DrawID: 1, Numbers: "5,12,23,34,45", Stars: "3,7", Cost: 2.50, PurchasedAt: time.Now()},
			}, nil
		},
	}

	rr := doTicketRequest(t, mock, "GET", "/tickets", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}

	var result []models.Ticket
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 ticket, got %d", len(result))
	}
}

func TestTicketHandler_Create_Success(t *testing.T) {
	mock := &mockTicketStore{
		createFn: func(drawID int, numbers, stars string, cost float64) (models.Ticket, error) {
			return models.Ticket{ID: 1, DrawID: drawID, Numbers: numbers, Stars: stars, Cost: cost, PurchasedAt: time.Now()}, nil
		},
	}

	rr := doTicketRequest(t, mock, "POST", "/tickets",
		`{"draw_id": 1, "numbers": "5,12,23,34,45", "stars": "3,7", "cost": 2.50}`)

	if rr.Code != http.StatusCreated {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusCreated)
	}
}

func TestTicketHandler_Create_MissingDrawID(t *testing.T) {
	mock := &mockTicketStore{}
	rr := doTicketRequest(t, mock, "POST", "/tickets",
		`{"numbers": "5,12,23,34,45", "cost": 2.50}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_Create_MissingNumbers(t *testing.T) {
	mock := &mockTicketStore{}
	rr := doTicketRequest(t, mock, "POST", "/tickets",
		`{"draw_id": 1, "cost": 2.50}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_Create_InvalidCost(t *testing.T) {
	mock := &mockTicketStore{}
	rr := doTicketRequest(t, mock, "POST", "/tickets",
		`{"draw_id": 1, "numbers": "5,12,23,34,45", "cost": 0}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestTicketHandler_UpdatePrize_Success(t *testing.T) {
	mock := &mockTicketStore{
		updatePrizeFn: func(id int, prizeTier *string, prizeAmount *float64, matchedNumbers *int, matchedStars *int) (models.Ticket, error) {
			return models.Ticket{
				ID: id, DrawID: 1, Numbers: "5,12,23,34,45", Stars: "3,7", Cost: 2.50,
				PrizeTier: prizeTier, PrizeAmount: prizeAmount, MatchedNumbers: matchedNumbers, MatchedStars: matchedStars,
			}, nil
		},
	}

	rr := doTicketRequest(t, mock, "PUT", "/tickets/1/prize",
		`{"prize_tier": "5th", "prize_amount": 50.0, "matched_numbers": 3, "matched_stars": 1}`)

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}

	var result models.Ticket
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if *result.PrizeTier != "5th" {
		t.Errorf("prize_tier = %v, want 5th", *result.PrizeTier)
	}
}

func TestTicketHandler_Delete_Success(t *testing.T) {
	mock := &mockTicketStore{
		deleteFn: func(id int) error {
			return nil
		},
	}

	rr := doTicketRequest(t, mock, "DELETE", "/tickets/1", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestTicketHandler_Delete_NotFound(t *testing.T) {
	mock := &mockTicketStore{
		deleteFn: func(id int) error {
			return errors.New("ticket 99 not found")
		},
	}

	rr := doTicketRequest(t, mock, "DELETE", "/tickets/99", "")

	if rr.Code != http.StatusNotFound {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusNotFound)
	}
}
