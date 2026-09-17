package handlers

import (
	"encoding/json"
	"net/http"

	"lottery-pool-manager/store"
	"lottery-pool-manager/utils"
)

type DashboardHandler struct {
	store *store.PostgreSQLStore
}

func NewDashboardHandler(s *store.PostgreSQLStore) *DashboardHandler {
	return &DashboardHandler{store: s}
}

func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetDashboardStats()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to fetch dashboard stats")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
