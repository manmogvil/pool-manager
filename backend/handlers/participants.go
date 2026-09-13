package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"lottery-pool-manager/models"
	"lottery-pool-manager/utils"
)

type ParticipantStore interface {
	GetAllParticipants() ([]models.Participant, error)
	GetParticipantByID(id int) (models.Participant, error)
	CreateParticipant(name, email string) (models.Participant, error)
	UpdateParticipant(id int, name, email string) (models.Participant, error)
	DeactivateParticipant(id int) (models.Participant, error)
	ActivateParticipant(id int) (models.Participant, error)
}

type ParticipantHandler struct {
	store ParticipantStore
}

func NewParticipantHandler(s ParticipantStore) *ParticipantHandler {
	return &ParticipantHandler{store: s}
}

func (h *ParticipantHandler) List(w http.ResponseWriter, r *http.Request) {
	participants, err := h.store.GetAllParticipants()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error fetching participants")
		return
	}
	utils.RespondJSON(w, http.StatusOK, participants)
}

func (h *ParticipantHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	p, err := h.store.GetParticipantByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, p)
}

type createParticipantRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (h *ParticipantHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, "email is required")
		return
	}

	p, err := h.store.CreateParticipant(req.Name, req.Email)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "error creating participant")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, p)
}

func (h *ParticipantHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	var req createParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, "name is required")
		return
	}
	if req.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, "email is required")
		return
	}

	p, err := h.store.UpdateParticipant(id, req.Name, req.Email)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, p)
}

func (h *ParticipantHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	p, err := h.store.DeactivateParticipant(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, p)
}

func (h *ParticipantHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ParseID(chi.URLParam(r, "id"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid participant id")
		return
	}

	p, err := h.store.ActivateParticipant(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, err.Error())
		return
	}

	utils.RespondJSON(w, http.StatusOK, p)
}
