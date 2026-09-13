package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"lottery-pool-manager/models"
)

func TestGameHandler_List_Success(t *testing.T) {
	mock := &mockLotteryGameStore{
		getAllFn: func() ([]models.LotteryGame, error) {
			return []models.LotteryGame{
				{ID: 1, Name: "EuroMillones", DrawDays: "Tuesday, Friday", TicketPrice: 2.50, Active: true},
			}, nil
		},
	}

	rr := doGameRequest(t, mock, "GET", "/games", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}

	var result []models.LotteryGame
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 game, got %d", len(result))
	}
}

func TestGameHandler_Create_Success(t *testing.T) {
	mock := &mockLotteryGameStore{
		createFn: func(name, drawDays string, ticketPrice float64) (models.LotteryGame, error) {
			return models.LotteryGame{ID: 1, Name: name, DrawDays: drawDays, TicketPrice: ticketPrice, Active: true}, nil
		},
	}

	rr := doGameRequest(t, mock, "POST", "/games",
		`{"name": "EuroMillones", "draw_days": "Tuesday, Friday", "ticket_price": 2.50}`)

	if rr.Code != http.StatusCreated {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusCreated)
	}
}

func TestGameHandler_Create_MissingName(t *testing.T) {
	mock := &mockLotteryGameStore{}
	rr := doGameRequest(t, mock, "POST", "/games",
		`{"draw_days": "Tuesday", "ticket_price": 2.50}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestGameHandler_Create_MissingDrawDays(t *testing.T) {
	mock := &mockLotteryGameStore{}
	rr := doGameRequest(t, mock, "POST", "/games",
		`{"name": "EuroMillones", "ticket_price": 2.50}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestGameHandler_Create_InvalidPrice(t *testing.T) {
	mock := &mockLotteryGameStore{}
	rr := doGameRequest(t, mock, "POST", "/games",
		`{"name": "EuroMillones", "draw_days": "Tuesday", "ticket_price": 0}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestGameHandler_Update_Success(t *testing.T) {
	mock := &mockLotteryGameStore{
		updateFn: func(id int, name, drawDays string, ticketPrice float64, active bool) (models.LotteryGame, error) {
			return models.LotteryGame{ID: id, Name: name, DrawDays: drawDays, TicketPrice: ticketPrice, Active: active}, nil
		},
	}

	rr := doGameRequest(t, mock, "PUT", "/games/1",
		`{"name": "Updated Game", "draw_days": "Monday", "ticket_price": 3.00, "active": true}`)

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestGameHandler_Delete_Success(t *testing.T) {
	mock := &mockLotteryGameStore{
		deleteFn: func(id int) error {
			return nil
		},
	}

	rr := doGameRequest(t, mock, "DELETE", "/games/1", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestGameHandler_Delete_HasTickets(t *testing.T) {
	mock := &mockLotteryGameStore{
		deleteFn: func(id int) error {
			return errors.New("cannot delete game 1: tickets exist for this game")
		},
	}

	rr := doGameRequest(t, mock, "DELETE", "/games/1", "")

	if rr.Code != http.StatusNotFound {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusNotFound)
	}
}
