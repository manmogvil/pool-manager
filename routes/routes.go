package routes

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"lottery-pool-manager/handlers"
	"lottery-pool-manager/store"
	"time"
)

func Setup(s *store.PostgreSQLStore) *chi.Mux {
	participantHandler := handlers.NewParticipantHandler(s)
	contributionHandler := handlers.NewContributionHandler(s)
	lotteryGameHandler := handlers.NewLotteryGameHandler(s)
	drawHandler := handlers.NewDrawHandler(s)
	ticketHandler := handlers.NewTicketHandler(s)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Route("/participants", func(r chi.Router) {
		r.Get("/", participantHandler.List)
		r.Get("/{id}", participantHandler.GetByID)
		r.Post("/", participantHandler.Create)
		r.Put("/{id}", participantHandler.Update)
		r.Delete("/{id}", participantHandler.Deactivate)
	})

	r.Route("/contributions", func(r chi.Router) {
		r.Get("/", contributionHandler.List)
		r.Get("/{id}", contributionHandler.GetByID)
		r.Post("/", contributionHandler.Create)
		r.Put("/{id}", contributionHandler.Update)
		r.Delete("/{id}", contributionHandler.Delete)
		r.Get("/participant/{id}", contributionHandler.GetByParticipant)
		r.Get("/period", contributionHandler.GetByPeriod)
	})

	r.Route("/games", func(r chi.Router) {
		r.Get("/", lotteryGameHandler.List)
		r.Get("/{id}", lotteryGameHandler.GetByID)
		r.Post("/", lotteryGameHandler.Create)
		r.Put("/{id}", lotteryGameHandler.Update)
		r.Delete("/{id}", lotteryGameHandler.Delete)
	})

	r.Route("/draws", func(r chi.Router) {
		r.Get("/", drawHandler.List)
		r.Get("/{id}", drawHandler.GetByID)
		r.Post("/", drawHandler.Create)
		r.Put("/{id}/results", drawHandler.UpdateResults)
		r.Put("/{id}/process", drawHandler.MarkAsProcessed)
		r.Delete("/{id}", drawHandler.Delete)
		r.Get("/game/{id}", drawHandler.GetByGame)
		r.Get("/pending", drawHandler.GetPending)
	})

	r.Route("/tickets", func(r chi.Router) {
		r.Get("/", ticketHandler.List)
		r.Get("/{id}", ticketHandler.GetByID)
		r.Post("/", ticketHandler.Create)
		r.Put("/{id}/prize", ticketHandler.UpdatePrize)
		r.Delete("/{id}", ticketHandler.Delete)
		r.Get("/draw/{id}", ticketHandler.GetByDraw)
	})

	return r
}
