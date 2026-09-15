package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"lottery-pool-manager/handlers"
)

func doRequest(t *testing.T, handler http.Handler, method, url, body string) *httptest.ResponseRecorder {
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

func doGameRequest(t *testing.T, mock *mockLotteryGameStore, method, url, body string) *httptest.ResponseRecorder {
	t.Helper()
	h := handlers.NewLotteryGameHandler(mock)
	return doRequest(t, setupGameRouter(h), method, url, body)
}

func doContributionRequest(t *testing.T, mock *mockContributionStore, method, url, body string) *httptest.ResponseRecorder {
	t.Helper()
	h := handlers.NewContributionHandler(mock)
	return doRequest(t, setupContributionRouter(h), method, url, body)
}

func doDrawRequest(t *testing.T, mock *mockDrawStore, method, url, body string) *httptest.ResponseRecorder {
	t.Helper()
	h := handlers.NewDrawHandler(mock)
	return doRequest(t, setupDrawRouter(h), method, url, body)
}

func doTicketRequest(t *testing.T, mock *mockTicketStore, method, url, body string) *httptest.ResponseRecorder {
	t.Helper()
	h := handlers.NewTicketHandler(mock)
	return doRequest(t, setupTicketRouter(h), method, url, body)
}

// Routers

func setupGameRouter(h *handlers.LotteryGameHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/games", h.List)
	r.Get("/games/{id}", h.GetByID)
	r.Post("/games", h.Create)
	r.Put("/games/{id}", h.Update)
	r.Delete("/games/{id}", h.Delete)
	return r
}

func setupContributionRouter(h *handlers.ContributionHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/contributions", h.List)
	r.Get("/contributions/{id}", h.GetByID)
	r.Post("/contributions", h.Create)
	r.Put("/contributions/{id}", h.Update)
	r.Delete("/contributions/{id}", h.Delete)
	return r
}

func setupDrawRouter(h *handlers.DrawHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/draws", h.List)
	r.Get("/draws/{id}", h.GetByID)
	r.Post("/draws", h.Create)
	r.Put("/draws/{id}/results", h.UpdateResults)
	r.Put("/draws/{id}/process", h.MarkAsProcessed)
	r.Delete("/draws/{id}", h.Delete)
	return r
}

func setupTicketRouter(h *handlers.TicketHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Get("/tickets", h.List)
	r.Get("/tickets/{id}", h.GetByID)
	r.Post("/tickets", h.Create)
	r.Put("/tickets/{id}/prize", h.UpdatePrize)
	r.Delete("/tickets/{id}", h.Delete)
	return r
}
