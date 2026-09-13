package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"lottery-pool-manager/models"
)

func TestDrawHandler_List_Success(t *testing.T) {
	mock := &mockDrawStore{
		getAllFn: func() ([]models.Draw, error) {
			return []models.Draw{
				{ID: 1, DrawDate: time.Now(), Processed: false},
			}, nil
		},
	}

	rr := doDrawRequest(t, mock, "GET", "/draws", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}

	var result []models.Draw
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 draw, got %d", len(result))
	}
}

func TestDrawHandler_Create_Success(t *testing.T) {
	mock := &mockDrawStore{
		createFn: func(gameID int, drawDate time.Time) (models.Draw, error) {
			return models.Draw{ID: 1, DrawDate: drawDate, Processed: false}, nil
		},
	}

	rr := doDrawRequest(t, mock, "POST", "/draws",
		`{"game_id": 1, "draw_date": "2024-01-15"}`)

	if rr.Code != http.StatusCreated {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusCreated)
	}
}

func TestDrawHandler_Create_MissingGameID(t *testing.T) {
	mock := &mockDrawStore{}
	rr := doDrawRequest(t, mock, "POST", "/draws",
		`{"draw_date": "2024-01-15"}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestDrawHandler_Create_MissingDate(t *testing.T) {
	mock := &mockDrawStore{}
	rr := doDrawRequest(t, mock, "POST", "/draws",
		`{"game_id": 1}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestDrawHandler_Create_InvalidDate(t *testing.T) {
	mock := &mockDrawStore{}
	rr := doDrawRequest(t, mock, "POST", "/draws",
		`{"game_id": 1, "draw_date": "15/01/2024"}`)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusBadRequest)
	}
}

func TestDrawHandler_Create_Duplicate(t *testing.T) {
	mock := &mockDrawStore{
		createFn: func(gameID int, drawDate time.Time) (models.Draw, error) {
			return models.Draw{}, errors.New("a draw already exists for this game on 2024-01-15")
		},
	}

	rr := doDrawRequest(t, mock, "POST", "/draws",
		`{"game_id": 1, "draw_date": "2024-01-15"}`)

	if rr.Code != http.StatusConflict {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusConflict)
	}
}

func TestDrawHandler_UpdateResults_Success(t *testing.T) {
	mock := &mockDrawStore{
		updateResultsFn: func(id int, resultNumbers *string, resultStars *string) (models.Draw, error) {
			return models.Draw{ID: id, DrawDate: time.Now(), ResultNumbers: resultNumbers, ResultStars: resultStars, Processed: false}, nil
		},
	}

	rr := doDrawRequest(t, mock, "PUT", "/draws/1/results",
		`{"result_numbers": "5,12,23,34,45", "result_stars": "3,7"}`)

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}

	var result models.Draw
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if *result.ResultNumbers != "5,12,23,34,45" {
		t.Errorf("result_numbers = %v, want 5,12,23,34,45", *result.ResultNumbers)
	}
}

func TestDrawHandler_MarkAsProcessed_Success(t *testing.T) {
	mock := &mockDrawStore{
		markProcessedFn: func(id int) error {
			return nil
		},
	}

	rr := doDrawRequest(t, mock, "PUT", "/draws/1/process", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestDrawHandler_Delete_Success(t *testing.T) {
	mock := &mockDrawStore{
		deleteFn: func(id int) error {
			return nil
		},
	}

	rr := doDrawRequest(t, mock, "DELETE", "/draws/1", "")

	if rr.Code != http.StatusOK {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusOK)
	}
}

func TestDrawHandler_Delete_HasTickets(t *testing.T) {
	mock := &mockDrawStore{
		deleteFn: func(id int) error {
			return errors.New("cannot delete draw 1: tickets exist for this draw")
		},
	}

	rr := doDrawRequest(t, mock, "DELETE", "/draws/1", "")

	if rr.Code != http.StatusNotFound {
		t.Errorf("status code = %v, want %v", rr.Code, http.StatusNotFound)
	}
}
