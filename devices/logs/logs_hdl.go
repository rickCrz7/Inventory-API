package logs

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rickCrz7/Inventory-API/utils"
)

type Handler struct {
	svc *Service
	// atz *authz.Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{
		svc: svc,
		// atz: atz,
	}
}

func (h *Handler) GetLogs(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	logs, err := h.svc.GetLogs(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(logs)
}

func (h *Handler) CreateLog(w http.ResponseWriter, r *http.Request) {
	var logEntry utils.DeviceLog
	if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.CreateLog(r.Context(), &logEntry); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(logEntry)
}

func (h *Handler) DeleteLog(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.svc.DeleteLog(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
