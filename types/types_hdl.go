package types

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

func (h *Handler) GetType(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	typ, err := h.svc.GetType(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(typ)
}

func (h *Handler) GetTypes(w http.ResponseWriter, r *http.Request) {
	types, err := h.svc.GetTypes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(types)
}

func (h *Handler) CreateType(w http.ResponseWriter, r *http.Request) {
	var typ utils.Type
	if err := json.NewDecoder(r.Body).Decode(&typ); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.CreateType(r.Context(), &typ); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(typ)
}

func (h *Handler) UpdateType(w http.ResponseWriter, r *http.Request) {
	var typ utils.Type
	if err := json.NewDecoder(r.Body).Decode(&typ); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.UpdateType(r.Context(), &typ); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(typ)
}

func (h *Handler) DeleteType(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.svc.DeleteType(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
