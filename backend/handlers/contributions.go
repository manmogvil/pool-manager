package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"lottery-pool-manager/models"
	"lottery-pool-manager/utils"
)

type ContributionStore interface {
	CreateContribution(c *models.Contribution) (models.Contribution, error)
	GetAllContributions() ([]models.Contribution, error)
	GetContributionByID(id int) (models.Contribution, error)
	GetContributionsByParticipant(participantID int) ([]models.Contribution, error)
	GetContributionsByPeriod(month, year int) ([]models.Contribution, error)
	UpdateContribution(id int, c *models.Contribution) (models.Contribution, error)
	DeleteContribution(id int) error
}

type ContributionHandler struct {
	store ContributionStore
}

func NewContributionHandler(s ContributionStore) *ContributionHandler {
	return &ContributionHandler{store: s}
}

func (h *ContributionHandler) List(w http.ResponseWriter, r *http.Request) {
	contributions, err := h.store.GetAllContributions()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error fetching contributions")
		return
	}
	utils.RespondJSON(w, http.StatusOK, contributions)
}

func (h *ContributionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid contribution id")
		return
	}

	c, err := h.store.GetContributionByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, c)
}

func (h *ContributionHandler) GetByParticipant(w http.ResponseWriter, r *http.Request) {
	participantID, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	contributions, err := h.store.GetContributionsByParticipant(participantID)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error fetching contributions")
		return
	}

	utils.RespondJSON(w, http.StatusOK, contributions)
}

func (h *ContributionHandler) GetByPeriod(w http.ResponseWriter, r *http.Request) {
	month, err := utils.ParseID(r.URL.Query().Get("month"))
	if err != nil || month < 1 || month > 12 {
		utils.RespondError(w, http.StatusBadRequest, "invalid month")
		return
	}
	year, err := utils.ParseID(r.URL.Query().Get("year"))
	if err != nil || year < 2020 || year > 2100 {
		utils.RespondError(w, http.StatusBadRequest, "invalid year")
		return
	}

	contributions, err := h.store.GetContributionsByPeriod(month, year)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error fetching contributions")
		return
	}

	utils.RespondJSON(w, http.StatusOK, contributions)
}

type createContributionRequest struct {
	ParticipantID int     `json:"participant_id"`
	GameID        int     `json:"game_id"`
	Month         int     `json:"month"`
	Year          int     `json:"year"`
	Amount        float64 `json:"amount"`
	Paid          bool    `json:"paid"`
	PaymentDate   *string `json:"payment_date"`
	PaymentMethod string  `json:"payment_method"`
	Comments      string  `json:"comments"`
}

func validateContributionRequest(req *createContributionRequest) error {
	if req.ParticipantID == 0 {
		return &utils.ValidationError{Field: "participant_id", Message: "is required"}
	}
	if req.GameID == 0 {
		return &utils.ValidationError{Field: "game_id", Message: "is required"}
	}
	if req.Month < 1 || req.Month > 12 {
		return &utils.ValidationError{Field: "month", Message: "must be between 1 and 12"}
	}
	if req.Year < 2020 {
		return &utils.ValidationError{Field: "year", Message: "must be 2020 or later"}
	}
	if req.Amount <= 0 {
		return &utils.ValidationError{Field: "amount", Message: "must be greater than 0"}
	}
	if req.PaymentMethod == "" {
		return &utils.ValidationError{Field: "payment_method", Message: "is required"}
	}
	if req.PaymentMethod != "CASH" && req.PaymentMethod != "BIZUM" {
		return &utils.ValidationError{Field: "payment_method", Message: "must be 'CASH' or 'BIZUM'"}
	}
	return nil
}

func buildContributionFromRequest(req *createContributionRequest) (*models.Contribution, error) {
	contribution := &models.Contribution{
		ParticipantID: req.ParticipantID,
		GameID:        req.GameID,
		Month:         req.Month,
		Year:          req.Year,
		Amount:        req.Amount,
		Paid:          req.Paid,
		PaymentMethod: req.PaymentMethod,
		Comments:      req.Comments,
	}

	if req.PaymentDate != nil {
		date, err := utils.ParseDate(*req.PaymentDate)
		if err != nil {
			return nil, &utils.ValidationError{Field: "payment_date", Message: "invalid format, use YYYY-MM-DD"}
		}
		contribution.PaymentDate = &date
	}

	return contribution, nil
}

func (h *ContributionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createContributionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validateContributionRequest(&req); err != nil {
		utils.RespondValidationError(w, err.(*utils.ValidationError))
		return
	}

	c, err := buildContributionFromRequest(&req)
	if err != nil {
		utils.RespondValidationError(w, err.(*utils.ValidationError))
		return
	}

	result, err := h.store.CreateContribution(c)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error creating contribution")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, result)
}

func (h *ContributionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid contribution id")
		return
	}

	var req createContributionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validateContributionRequest(&req); err != nil {
		utils.RespondValidationError(w, err.(*utils.ValidationError))
		return
	}

	c, err := buildContributionFromRequest(&req)
	if err != nil {
		utils.RespondValidationError(w, err.(*utils.ValidationError))
		return
	}

	result, err := h.store.UpdateContribution(id, c)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, result)
}

func (h *ContributionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid contribution id")
		return
	}

	if err := h.store.DeleteContribution(id); err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "contribution deleted"})
}
