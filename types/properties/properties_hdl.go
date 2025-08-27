package properties

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

func (h *Handler) GetProperties(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["type_id"]
	properties, err := h.svc.GetProperties(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(properties)
}

func (h *Handler) CreateProperty(w http.ResponseWriter, r *http.Request) {
	var property utils.TypeProperty
	if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.CreateProperty(r.Context(), &property); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(property)
}

func (h *Handler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	var property utils.TypeProperty
	if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.UpdateProperty(r.Context(), &property); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(property)
}

func (h *Handler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.svc.DeleteProperty(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
