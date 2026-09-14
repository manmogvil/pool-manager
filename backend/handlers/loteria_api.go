package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"lottery-pool-manager/models"
	"lottery-pool-manager/services"
	"lottery-pool-manager/utils"
)

type LoteriaAPIStore interface {
	GetLotteryGameByName(name string) (models.LotteryGame, error)
	CreateDraw(gameID int, drawDate time.Time) (models.Draw, error)
	GetDrawsByGame(gameID int) ([]models.Draw, error)
	UpdateDrawResults(id int, resultNumbers *string, resultStars *string) (models.Draw, error)
	UpdateDrawDrawIDAPI(id int, drawIDAPI string) (models.Draw, error)
}

type LoteriaAPIHandler struct {
	store LoteriaAPIStore
	api   services.LoteriaAPI
}

func NewLoteriaAPIHandler(s LoteriaAPIStore) *LoteriaAPIHandler {
	return &LoteriaAPIHandler{
		store: s,
		api:   services.NewLoteriaAPIClient(),
	}
}

func NewLoteriaAPIHandlerWithAPI(s LoteriaAPIStore, api services.LoteriaAPI) *LoteriaAPIHandler {
	return &LoteriaAPIHandler{
		store: s,
		api:   api,
	}
}

type FetchResultsRequest struct {
	GameSlug string `json:"game_slug"`
	From     string `json:"from"`
	To       string `json:"to"`
}

type FetchResultsResponse struct {
	Draws   []models.Draw `json:"draws"`
	Message string        `json:"message"`
}

func (h *LoteriaAPIHandler) FetchResults(w http.ResponseWriter, r *http.Request) {
	var req FetchResultsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.GameSlug == "" {
		utils.RespondError(w, http.StatusBadRequest, "game_slug is required")
		return
	}

	if req.From == "" || req.To == "" {
		utils.RespondError(w, http.StatusBadRequest, "from and to dates are required (YYYY-MM-DD)")
		return
	}

	result, err := h.api.GetResults(req.GameSlug, req.From, req.To)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching from Loteria API: "+err.Error())
		return
	}

	gameName := services.MapGameSlugToName(req.GameSlug)
	game, err := h.store.GetLotteryGameByName(gameName)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Game not found: "+gameName)
		return
	}

	existingDraws, _ := h.store.GetDrawsByGame(game.ID)
	existingByDate := make(map[string]models.Draw)
	for _, d := range existingDraws {
		existingByDate[d.DrawDate.Format("2006-01-02")] = d
	}

	var updatedDraws []models.Draw

	for _, item := range result.Data {
		drawDate, err := time.Parse("2006-01-02", item.DrawDate)
		if err != nil {
			continue
		}

		resultNumbers := services.FormatCombination(item.Combination)
		resultStars := services.FormatStars(item.ResultData.Estrellas)
		dateKey := drawDate.Format("2006-01-02")

		if existing, ok := existingByDate[dateKey]; ok {
			updated, err := h.store.UpdateDrawResults(existing.ID, &resultNumbers, &resultStars)
			if err != nil {
				continue
			}
			if item.DrawId != "" {
				updated, err = h.store.UpdateDrawDrawIDAPI(existing.ID, item.DrawId)
				if err != nil {
					continue
				}
			}
			updatedDraws = append(updatedDraws, updated)
		}
	}

	utils.RespondJSON(w, http.StatusOK, FetchResultsResponse{
		Draws:   updatedDraws,
		Message: fmt.Sprintf("Processed %d draws from Loteria API", len(updatedDraws)),
	})
}
