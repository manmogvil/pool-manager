package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"lottery-pool-manager/models"
	"lottery-pool-manager/utils"
)

type DrawStore interface {
	CreateDraw(gameID *int, drawDate time.Time) (models.Draw, error)
	GetAllDraws() ([]models.Draw, error)
	GetDrawByID(id int) (models.Draw, error)
	GetDrawsByGame(gameID int) ([]models.Draw, error)
	GetPendingDraws() ([]models.Draw, error)
	UpdateDrawResults(id int, resultNumbers *string, resultStars *string) (models.Draw, error)
	MarkDrawAsProcessed(id int) error
	DeleteDraw(id int) error
}

type DrawHandler struct {
	store DrawStore
}

func NewDrawHandler(s DrawStore) *DrawHandler {
	return &DrawHandler{store: s}
}

func (h *DrawHandler) List(w http.ResponseWriter, r *http.Request) {
	draws, err := h.store.GetAllDraws()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error fetching draws")
		return
	}
	utils.RespondJSON(w, http.StatusOK, draws)
}

func (h *DrawHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid draw id")
		return
	}

	d, err := h.store.GetDrawByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, d)
}

func (h *DrawHandler) GetByGame(w http.ResponseWriter, r *http.Request) {
	gameID, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid game id")
		return
	}

	draws, err := h.store.GetDrawsByGame(gameID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error fetching draws")
		return
	}

	utils.RespondJSON(w, http.StatusOK, draws)
}

func (h *DrawHandler) GetPending(w http.ResponseWriter, r *http.Request) {
	draws, err := h.store.GetPendingDraws()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error fetching pending draws")
		return
	}
	utils.RespondJSON(w, http.StatusOK, draws)
}

type createDrawRequest struct {
	GameID   *int   `json:"game_id"`
	DrawDate string `json:"draw_date"`
}

func (h *DrawHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createDrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DrawDate == "" {
		utils.RespondError(w, http.StatusBadRequest, "draw_date is required")
		return
	}

	drawDate, err := utils.ParseDate(req.DrawDate)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid draw_date format, use YYYY-MM-DD")
		return
	}

	d, err := h.store.CreateDraw(req.GameID, drawDate)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error creating draw")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, d)
}

type updateDrawResultsRequest struct {
	ResultNumbers *string `json:"result_numbers"`
	ResultStars   *string `json:"result_stars"`
}

func (h *DrawHandler) UpdateResults(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid draw id")
		return
	}

	var req updateDrawResultsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ResultNumbers == nil || *req.ResultNumbers == "" {
		utils.RespondError(w, http.StatusBadRequest, "result_numbers is required")
		return
	}

	d, err := h.store.UpdateDrawResults(id, req.ResultNumbers, req.ResultStars)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, d)
}

func (h *DrawHandler) MarkAsProcessed(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid draw id")
		return
	}

	if err := h.store.MarkDrawAsProcessed(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "draw marked as processed"})
}

func (h *DrawHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid draw id")
		return
	}

	if err := h.store.DeleteDraw(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "draw deleted"})
}
