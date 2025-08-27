package owners

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

func (h *Handler) GetOwner(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	owner, err := h.svc.GetOwner(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(owner)
}

func (h *Handler) GetOwnerByCampusID(w http.ResponseWriter, r *http.Request) {
	campusID := mux.Vars(r)["campusID"]
	owner, err := h.svc.GetOwnerByCampusID(r.Context(), campusID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(owner)
}

func (h *Handler) GetOwnerByEmail(w http.ResponseWriter, r *http.Request) {
	email := mux.Vars(r)["email"]
	owner, err := h.svc.GetOwnerByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(owner)
}

func (h *Handler) GetOwners(w http.ResponseWriter, r *http.Request) {
	owners, err := h.svc.GetOwners(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(owners)
}

func (h *Handler) CreateOwner(w http.ResponseWriter, r *http.Request) {
	var owner utils.Owner
	if err := json.NewDecoder(r.Body).Decode(&owner); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.CreateOwner(r.Context(), &owner); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(owner)
}

func (h *Handler) UpdateOwner(w http.ResponseWriter, r *http.Request) {
	var owner utils.Owner
	if err := json.NewDecoder(r.Body).Decode(&owner); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.svc.UpdateOwner(r.Context(), &owner); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(owner)
}

func (h *Handler) DeleteOwner(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.svc.DeleteOwner(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
