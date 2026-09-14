package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"lottery-pool-manager/models"
	"lottery-pool-manager/services"
	"lottery-pool-manager/utils"
)

type CheckTicketStore interface {
	GetTicketByID(id int) (models.Ticket, error)
	GetDrawByID(id int) (models.Draw, error)
	GetLotteryGameByID(id int) (models.LotteryGame, error)
	UpdateTicketPrize(id int, prizeTier *string, prizeAmount *float64, matchedNumbers *int, matchedStars *int) (models.Ticket, error)
}

type CheckTicketHandler struct {
	store CheckTicketStore
	api   services.LoteriaAPI
}

func NewCheckTicketHandler(s CheckTicketStore) *CheckTicketHandler {
	return &CheckTicketHandler{
		store: s,
		api:   services.NewLoteriaAPIClient(),
	}
}

func NewCheckTicketHandlerWithAPI(s CheckTicketStore, api services.LoteriaAPI) *CheckTicketHandler {
	return &CheckTicketHandler{
		store: s,
		api:   api,
	}
}

type CheckTicketRequest struct {
	TicketID int `json:"ticket_id"`
}

type CheckTicketResponse struct {
	Ticket         models.Ticket `json:"ticket"`
	HasPrize       bool          `json:"has_prize"`
	Category       string        `json:"category"`
	PrizeAmount    string        `json:"prize_amount"`
	MatchedNumbers int           `json:"matched_numbers"`
	MatchedStars   int           `json:"matched_stars"`
	Message        string        `json:"message"`
}

func (h *CheckTicketHandler) CheckTicket(w http.ResponseWriter, r *http.Request) {
	var req CheckTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.TicketID == 0 {
		utils.RespondError(w, http.StatusBadRequest, "ticket_id is required")
		return
	}

	ticket, err := h.store.GetTicketByID(req.TicketID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Ticket not found")
		return
	}

	draw, err := h.store.GetDrawByID(ticket.DrawID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Draw not found")
		return
	}

	game, err := h.store.GetLotteryGameByID(draw.GameID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Game not found")
		return
	}

	gameSlug := services.MapNameToGameSlug(game.Name)
	if gameSlug == "" {
		utils.RespondError(w, http.StatusBadRequest, "Game not supported for API check")
		return
	}

	numbers := ticket.Numbers
	stars := ""
	if ticket.Stars != nil {
		stars = *ticket.Stars
	}

	drawIDAPI := ""
	if draw.DrawIDAPI != nil {
		drawIDAPI = *draw.DrawIDAPI
	}

	result, err := h.api.CheckCombination(gameSlug, numbers, stars, drawIDAPI)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error checking combination: "+err.Error())
		return
	}

	var prizeTier *string
	var prizeAmount *float64
	matchedNumbers := result.Data.MainNumbersMatched
	matchedStars := result.Data.ExtraNumbersMatched

	if result.Data.IsWinner {
		tier := fmt.Sprintf("%d matched", matchedNumbers)
		prizeTier = &tier
	}

	_, err = h.store.UpdateTicketPrize(ticket.ID, prizeTier, prizeAmount, &matchedNumbers, &matchedStars)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating ticket: "+err.Error())
		return
	}

	ticket.PrizeTier = prizeTier
	ticket.MatchedNumbers = &matchedNumbers
	ticket.MatchedStars = &matchedStars

	utils.RespondJSON(w, http.StatusOK, CheckTicketResponse{
		Ticket:         ticket,
		HasPrize:       result.Data.IsWinner,
		Category:       "",
		PrizeAmount:    "",
		MatchedNumbers: matchedNumbers,
		MatchedStars:   matchedStars,
		Message:        "Ticket checked successfully",
	})
}
