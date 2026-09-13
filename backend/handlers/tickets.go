package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"lottery-pool-manager/models"
	"lottery-pool-manager/utils"
)

type TicketStore interface {
	CreateTicket(drawID int, numbers string, stars *string, cost float64) (models.Ticket, error)
	GetAllTickets() ([]models.Ticket, error)
	GetTicketByID(id int) (models.Ticket, error)
	GetTicketsByDraw(drawID int) ([]models.Ticket, error)
	UpdateTicketPrize(id int, prizeTier *string, prizeAmount *float64, matchedNumbers *int, matchedStars *int) (models.Ticket, error)
	DeleteTicket(id int) error
}

type TicketHandler struct {
	store TicketStore
}

func NewTicketHandler(s TicketStore) *TicketHandler {
	return &TicketHandler{store: s}
}

func (h *TicketHandler) List(w http.ResponseWriter, r *http.Request) {
	tickets, err := h.store.GetAllTickets()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error fetching tickets")
		return
	}
	utils.RespondJSON(w, http.StatusOK, tickets)
}

func (h *TicketHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	t, err := h.store.GetTicketByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, t)
}

func (h *TicketHandler) GetByDraw(w http.ResponseWriter, r *http.Request) {
	drawID, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid draw id")
		return
	}

	tickets, err := h.store.GetTicketsByDraw(drawID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error fetching tickets")
		return
	}

	utils.RespondJSON(w, http.StatusOK, tickets)
}

type createTicketRequest struct {
	DrawID  int     `json:"draw_id"`
	Numbers string  `json:"numbers"`
	Stars   *string `json:"stars"`
	Cost    float64 `json:"cost"`
}

func validateTicketRequest(req *createTicketRequest) error {
	if req.DrawID == 0 {
		return &utils.ValidationError{Field: "draw_id", Message: "is required"}
	}
	if req.Numbers == "" {
		return &utils.ValidationError{Field: "numbers", Message: "is required"}
	}
	if req.Cost <= 0 {
		return &utils.ValidationError{Field: "cost", Message: "must be greater than 0"}
	}
	return nil
}

func (h *TicketHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validateTicketRequest(&req); err != nil {
		utils.RespondValidationError(w, err.(*utils.ValidationError))
		return
	}

	t, err := h.store.CreateTicket(req.DrawID, req.Numbers, req.Stars, req.Cost)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error creating ticket")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, t)
}

type updateTicketPrizeRequest struct {
	PrizeTier      *string  `json:"prize_tier"`
	PrizeAmount    *float64 `json:"prize_amount"`
	MatchedNumbers *int     `json:"matched_numbers"`
	MatchedStars   *int     `json:"matched_stars"`
}

func (h *TicketHandler) UpdatePrize(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	var req updateTicketPrizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	t, err := h.store.UpdateTicketPrize(id, req.PrizeTier, req.PrizeAmount, req.MatchedNumbers, req.MatchedStars)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, t)
}

func (h *TicketHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid ticket id")
		return
	}

	if err := h.store.DeleteTicket(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "ticket deleted"})
}
