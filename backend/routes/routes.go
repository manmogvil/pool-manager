package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"lottery-pool-manager/handlers"
	appMiddleware "lottery-pool-manager/middleware"
	"lottery-pool-manager/services"
	"lottery-pool-manager/store"
	"lottery-pool-manager/utils"
	"time"
)

func Setup(s *store.PostgreSQLStore) *chi.Mux {
	contributionHandler := handlers.NewContributionHandler(s)
	lotteryGameHandler := handlers.NewLotteryGameHandler(s)
	drawHandler := handlers.NewDrawHandler(s)
	ticketHandler := handlers.NewTicketHandler(s)
	authHandler := handlers.NewAuthHandler(s)
	loteriaAPIHandler := handlers.NewLoteriaAPIHandler(s)
	checkTicketHandler := handlers.NewCheckTicketHandler(s)
	dashboardHandler := handlers.NewDashboardHandler(s)

	authService := services.NewAuthService()

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   strings.Split(utils.GetEnv("CORS_ORIGIN", "http://localhost:5173"), ","),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Auth(authService))
			r.Get("/me", authHandler.GetMe)
			r.Get("/users/names", authHandler.ListUserNames)
			r.Group(func(r chi.Router) {
				r.Use(appMiddleware.RequireAdmin)
				r.Get("/users", authHandler.ListUsers)
				r.Put("/users/{id}", authHandler.UpdateUser)
				r.Put("/users/{id}/activate", authHandler.ActivateUser)
				r.Put("/users/{id}/deactivate", authHandler.DeactivateUser)
			})
		})
	})

	r.Route("/contributions", func(r chi.Router) {
		r.Get("/", contributionHandler.List)
		r.Get("/{id}", contributionHandler.GetByID)
		r.Get("/user/{id}", contributionHandler.GetByUser)
		r.Get("/period", contributionHandler.GetByPeriod)
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Auth(authService))
			r.Post("/", contributionHandler.Create)
			r.Put("/{id}", contributionHandler.Update)
			r.Delete("/{id}", contributionHandler.Delete)
		})
	})

	r.Route("/games", func(r chi.Router) {
		r.Get("/", lotteryGameHandler.List)
		r.Get("/{id}", lotteryGameHandler.GetByID)
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Auth(authService))
			r.Use(appMiddleware.RequireAdmin)
			r.Post("/", lotteryGameHandler.Create)
			r.Put("/{id}", lotteryGameHandler.Update)
			r.Delete("/{id}", lotteryGameHandler.Delete)
		})
	})

	r.Route("/draws", func(r chi.Router) {
		r.Get("/", drawHandler.List)
		r.Get("/{id}", drawHandler.GetByID)
		r.Get("/game/{id}", drawHandler.GetByGame)
		r.Get("/pending", drawHandler.GetPending)
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Auth(authService))
			r.Post("/", drawHandler.Create)
			r.Put("/{id}/results", drawHandler.UpdateResults)
			r.Put("/{id}/process", drawHandler.MarkAsProcessed)
			r.Delete("/{id}", drawHandler.Delete)
		})
	})

	r.Route("/tickets", func(r chi.Router) {
		r.Get("/", ticketHandler.List)
		r.Get("/{id}", ticketHandler.GetByID)
		r.Get("/draw/{id}", ticketHandler.GetByDraw)
		r.Group(func(r chi.Router) {
			r.Use(appMiddleware.Auth(authService))
			r.Post("/", ticketHandler.Create)
			r.Put("/{id}/prize", ticketHandler.UpdatePrize)
			r.Delete("/{id}", ticketHandler.Delete)
		})
	})

	r.Route("/loteria-api", func(r chi.Router) {
		r.Use(appMiddleware.Auth(authService))
		r.Use(appMiddleware.RequireAdmin)
		r.Post("/fetch-results", loteriaAPIHandler.FetchResults)
	})

	r.Group(func(r chi.Router) {
		r.Use(appMiddleware.Auth(authService))
		r.Use(appMiddleware.RequireAdmin)
		r.Post("/check-ticket", checkTicketHandler.CheckTicket)
	})

	r.Route("/dashboard", func(r chi.Router) {
		r.Use(appMiddleware.Auth(authService))
		r.Get("/stats", dashboardHandler.GetStats)
	})

	return r
}
