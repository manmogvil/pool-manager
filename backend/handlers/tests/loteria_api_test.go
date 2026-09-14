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

func setupLoteriaAPIRouter(h *handlers.LoteriaAPIHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/loteria-api/fetch-results", h.FetchResults)
	return r
}

func doLoteriaAPIRequest(t *testing.T, mock *mockLoteriaAPIStore, api services.LoteriaAPI, method, url, body string) *httptest.ResponseRecorder {
	t.Helper()
	h := handlers.NewLoteriaAPIHandlerWithAPI(mock, api)
	return doRequest(t, setupLoteriaAPIRouter(h), method, url, body)
}

func TestFetchResults_MissingFields(t *testing.T) {
	mock := &mockLoteriaAPIStore{}
	api := &mockLoteriaAPI{}
	rr := doLoteriaAPIRequest(t, mock, api, "POST", "/loteria-api/fetch-results", `{}`)
	if rr.Code != 400 {
		t.Errorf("expected 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestFetchResults_MissingDates(t *testing.T) {
	mock := &mockLoteriaAPIStore{}
	api := &mockLoteriaAPI{}
	rr := doLoteriaAPIRequest(t, mock, api, "POST", "/loteria-api/fetch-results", `{"game_slug": "euromillones"}`)
	if rr.Code != 400 {
		t.Errorf("expected 400, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestFetchResults_GameNotFound(t *testing.T) {
	mock := &mockLoteriaAPIStore{
		getLotteryGameByNameFn: func(name string) (models.LotteryGame, error) {
			return models.LotteryGame{}, fmt.Errorf("game not found")
		},
	}
	api := &mockLoteriaAPI{
		getResultsFn: func(gameSlug string, fromDate string, toDate string) (*services.ResultsListResponse, error) {
			return &services.ResultsListResponse{
				Success: true,
				Data:    []services.ResultsListItem{},
			}, nil
		},
	}
	rr := doLoteriaAPIRequest(t, mock, api, "POST", "/loteria-api/fetch-results", `{"game_slug": "euromillones", "from": "2026-09-08", "to": "2026-09-14"}`)
	if rr.Code != 404 {
		t.Errorf("expected 404, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestFetchResults_Success(t *testing.T) {
	mock := &mockLoteriaAPIStore{
		getLotteryGameByNameFn: func(name string) (models.LotteryGame, error) {
			return models.LotteryGame{ID: 1, Name: "EuroMillones"}, nil
		},
		getDrawsByGameFn: func(gameID int) ([]models.Draw, error) {
			return []models.Draw{}, nil
		},
		createDrawFn: func(gameID int, drawDate time.Time) (models.Draw, error) {
			return models.Draw{ID: 10, GameID: gameID, DrawDate: drawDate}, nil
		},
		updateDrawResultsFn: func(id int, resultNumbers *string, resultStars *string) (models.Draw, error) {
			return models.Draw{ID: id, ResultNumbers: resultNumbers, ResultStars: resultStars}, nil
		},
		updateDrawIDAPIFn: func(id int, drawIDAPI string) (models.Draw, error) {
			return models.Draw{ID: id, DrawIDAPI: &drawIDAPI}, nil
		},
	}

	api := &mockLoteriaAPI{
		getResultsFn: func(gameSlug string, fromDate string, toDate string) (*services.ResultsListResponse, error) {
			return &services.ResultsListResponse{
				Success: true,
				Data: []services.ResultsListItem{
					{
						DrawDate: "2026-09-11",
						Combination: []int{5, 12, 23, 34, 45},
						ResultData: services.LoteriaAPIResult{Estrellas: []int{3, 7}},
						DrawId: "1322202073",
					},
				},
			}, nil
		},
	}

	rr := doLoteriaAPIRequest(t, mock, api, "POST", "/loteria-api/fetch-results", `{"game_slug": "euromillones", "from": "2026-09-08", "to": "2026-09-14"}`)
	if rr.Code != 200 {
		t.Errorf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
}
