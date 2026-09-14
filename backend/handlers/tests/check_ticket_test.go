package tests

import (
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"lottery-pool-manager/handlers"
	"lottery-pool-manager/models"
	"lottery-pool-manager/services"
)

func setupCheckTicketRouter(h *handlers.CheckTicketHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/check-ticket", h.CheckTicket)
	return r
}

func doCheckTicketRequest(t *testing.T, mock *mockCheckTicketStore, api services.LoteriaAPI, method, url, body string) *httptest.ResponseRecorder {
	t.Helper()
	h := handlers.NewCheckTicketHandlerWithAPI(mock, api)
	return doRequest(t, setupCheckTicketRouter(h), method, url, body)
}

func strPtr(s string) *string {
	return &s
}

func TestCheckTicket_Success(t *testing.T) {
	mock := &mockCheckTicketStore{
		getTicketByIDFn: func(id int) (models.Ticket, error) {
			return models.Ticket{
				ID:      1,
				DrawID:  1,
				Numbers: "5,12,23,34,45",
				Stars:   strPtr("3,7"),
				Cost:    2.50,
			}, nil
		},
		getDrawByIDFn: func(id int) (models.Draw, error) {
			return models.Draw{
				ID:       1,
				GameID:   1,
				DrawDate: time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC),
			}, nil
		},
		getLotteryGameByIDFn: func(id int) (models.LotteryGame, error) {
			return models.LotteryGame{
				ID:   1,
				Name: "EuroMillones",
			}, nil
		},
		updateTicketPrizeFn: func(id int, prizeTier *string, prizeAmount *float64, matchedNumbers *int, matchedStars *int) (models.Ticket, error) {
			return models.Ticket{
				ID:             1,
				DrawID:         1,
				Numbers:        "5,12,23,34,45",
				Stars:          strPtr("3,7"),
				Cost:           2.50,
				PrizeTier:      prizeTier,
				MatchedNumbers: matchedNumbers,
				MatchedStars:   matchedStars,
			}, nil
		},
	}

	api := &mockLoteriaAPI{
		checkCombinationFn: func(gameSlug string, numbers string, extraNumbers string, drawId string) (*services.CheckCombinationResponse, error) {
			return &services.CheckCombinationResponse{
				Success: true,
				Data: services.CheckCombinationData{
					IsWinner:            true,
					MainNumbersMatched:  5,
					ExtraNumbersMatched: 2,
				},
			}, nil
		},
	}

	rr := doCheckTicketRequest(t, mock, api, "POST", "/check-ticket", `{"ticket_id": 1}`)
	if rr.Code != 200 {
		t.Errorf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestCheckTicket_MissingTicketID(t *testing.T) {
	mock := &mockCheckTicketStore{}
	api := &mockLoteriaAPI{}
	rr := doCheckTicketRequest(t, mock, api, "POST", "/check-ticket", `{}`)
	if rr.Code != 400 {
		t.Errorf("expected 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestCheckTicket_TicketNotFound(t *testing.T) {
	mock := &mockCheckTicketStore{
		getTicketByIDFn: func(id int) (models.Ticket, error) {
			return models.Ticket{}, fmt.Errorf("ticket 999 not found")
		},
	}
	api := &mockLoteriaAPI{}
	rr := doCheckTicketRequest(t, mock, api, "POST", "/check-ticket", `{"ticket_id": 999}`)
	if rr.Code != 404 {
		t.Errorf("expected 404, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestCheckTicket_DrawNotFound(t *testing.T) {
	mock := &mockCheckTicketStore{
		getTicketByIDFn: func(id int) (models.Ticket, error) {
			return models.Ticket{ID: 1, DrawID: 1}, nil
		},
		getDrawByIDFn: func(id int) (models.Draw, error) {
			return models.Draw{}, fmt.Errorf("draw 1 not found")
		},
	}
	api := &mockLoteriaAPI{}
	rr := doCheckTicketRequest(t, mock, api, "POST", "/check-ticket", `{"ticket_id": 1}`)
	if rr.Code != 404 {
		t.Errorf("expected 404, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestCheckTicket_GameNotFound(t *testing.T) {
	mock := &mockCheckTicketStore{
		getTicketByIDFn: func(id int) (models.Ticket, error) {
			return models.Ticket{ID: 1, DrawID: 1}, nil
		},
		getDrawByIDFn: func(id int) (models.Draw, error) {
			return models.Draw{ID: 1, GameID: 1}, nil
		},
		getLotteryGameByIDFn: func(id int) (models.LotteryGame, error) {
			return models.LotteryGame{}, fmt.Errorf("game 1 not found")
		},
	}
	api := &mockLoteriaAPI{}
	rr := doCheckTicketRequest(t, mock, api, "POST", "/check-ticket", `{"ticket_id": 1}`)
	if rr.Code != 404 {
		t.Errorf("expected 404, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestCheckTicket_UnsupportedGame(t *testing.T) {
	mock := &mockCheckTicketStore{
		getTicketByIDFn: func(id int) (models.Ticket, error) {
			return models.Ticket{ID: 1, DrawID: 1}, nil
		},
		getDrawByIDFn: func(id int) (models.Draw, error) {
			return models.Draw{ID: 1, GameID: 1}, nil
		},
		getLotteryGameByIDFn: func(id int) (models.LotteryGame, error) {
			return models.LotteryGame{ID: 1, Name: "Unknown Game"}, nil
		},
	}
	api := &mockLoteriaAPI{}
	rr := doCheckTicketRequest(t, mock, api, "POST", "/check-ticket", `{"ticket_id": 1}`)
	if rr.Code != 400 {
		t.Errorf("expected 400, got %d: %s", rr.Code, rr.Body.String())
	}
}
