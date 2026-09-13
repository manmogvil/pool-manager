package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"lottery-pool-manager/models"
	"lottery-pool-manager/utils"
)

type LotteryGameStore interface {
	CreateLotteryGame(name, drawDays string, ticketPrice float64) (models.LotteryGame, error)
	GetAllLotteryGames() ([]models.LotteryGame, error)
	GetLotteryGameByID(id int) (models.LotteryGame, error)
	UpdateLotteryGame(id int, name, drawDays string, ticketPrice float64, active bool) (models.LotteryGame, error)
	DeleteLotteryGame(id int) error
}

type LotteryGameHandler struct {
	store LotteryGameStore
}

func NewLotteryGameHandler(s LotteryGameStore) *LotteryGameHandler {
	return &LotteryGameHandler{store: s}
}

func (h *LotteryGameHandler) List(w http.ResponseWriter, r *http.Request) {
	games, err := h.store.GetAllLotteryGames()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error fetching lottery games")
		return
	}
	utils.RespondJSON(w, http.StatusOK, games)
}

func (h *LotteryGameHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid game id")
		return
	}

	g, err := h.store.GetLotteryGameByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, g)
}

type createLotteryGameRequest struct {
	Name        string  `json:"name"`
	DrawDays    string  `json:"draw_days"`
	TicketPrice float64 `json:"ticket_price"`
}

func validateLotteryGameRequest(req *createLotteryGameRequest) error {
	if req.Name == "" {
		return &utils.ValidationError{Field: "name", Message: "is required"}
	}
	if req.DrawDays == "" {
		return &utils.ValidationError{Field: "draw_days", Message: "is required"}
	}
	if req.TicketPrice <= 0 {
		return &utils.ValidationError{Field: "ticket_price", Message: "must be greater than 0"}
	}
	return nil
}

func (h *LotteryGameHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createLotteryGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validateLotteryGameRequest(&req); err != nil {
		utils.RespondValidationError(w, err.(*utils.ValidationError))
		return
	}

	g, err := h.store.CreateLotteryGame(req.Name, req.DrawDays, req.TicketPrice)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error creating lottery game")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, g)
}

func (h *LotteryGameHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid game id")
		return
	}

	type updateRequest struct {
		Name        string  `json:"name"`
		DrawDays    string  `json:"draw_days"`
		TicketPrice float64 `json:"ticket_price"`
		Active      bool    `json:"active"`
	}

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.DrawDays == "" {
		utils.RespondError(w, http.StatusBadRequest, "draw_days is required")
		return
	}
	if req.TicketPrice <= 0 {
		utils.RespondError(w, http.StatusBadRequest, "ticket_price must be greater than 0")
		return
	}

	g, err := h.store.UpdateLotteryGame(id, req.Name, req.DrawDays, req.TicketPrice, req.Active)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, g)
}

func (h *LotteryGameHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid game id")
		return
	}

	if err := h.store.DeleteLotteryGame(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "lottery game deleted"})
}
