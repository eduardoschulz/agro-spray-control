package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/eduardoschulz/agro-spray-control/backend/internal/repositories"
	"github.com/gorilla/mux"
)

type PulverizacaoAreaHandler struct {
	repo *repositories.PulverizacaoAreaRepo
}

func NewPulverizacaoAreaHandler(repo *repositories.PulverizacaoAreaRepo) *PulverizacaoAreaHandler {
	return &PulverizacaoAreaHandler{repo: repo}
}

// AssociateArea associa uma área a uma pulverização
func (h *PulverizacaoAreaHandler) AssociateArea(w http.ResponseWriter, r *http.Request) {
	var pa models.PulverizacaoArea
	if err := json.NewDecoder(r.Body).Decode(&pa); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	if err := h.repo.Associate(r.Context(), &pa); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pa)
}

// GetAreasByPulverizacao obtém as áreas de uma pulverização
func (h *PulverizacaoAreaHandler) GetAreasByPulverizacao(w http.ResponseWriter, r *http.Request) {
	codPulv := mux.Vars(r)["codPulv"]
	if codPulv == "" {
		http.Error(w, "Código de pulverização é obrigatório", http.StatusBadRequest)
		return
	}

	areas, err := h.repo.GetAreasByPulverizacao(r.Context(), codPulv)
	if err != nil {
		http.Error(w, "Erro ao buscar áreas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(areas)
}

// GetPulverizacoesByArea obtém as pulverizações de uma área
func (h *PulverizacaoAreaHandler) GetPulverizacoesByArea(w http.ResponseWriter, r *http.Request) {
	codArea := mux.Vars(r)["codArea"]
	if codArea == "" {
		http.Error(w, "Código de área é obrigatório", http.StatusBadRequest)
		return
	}

	pulverizacoes, err := h.repo.GetPulverizacoesByArea(r.Context(), codArea)
	if err != nil {
		http.Error(w, "Erro ao buscar pulverizações", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pulverizacoes)
}

// RemoveAssociation remove a associação entre uma área e uma pulverização
func (h *PulverizacaoAreaHandler) RemoveAssociation(w http.ResponseWriter, r *http.Request) {
	codPulv := mux.Vars(r)["codPulv"]
	codArea := mux.Vars(r)["codArea"]

	if codPulv == "" || codArea == "" {
		http.Error(w, "Códigos de pulverização e área são obrigatórios", http.StatusBadRequest)
		return
	}

	if err := h.repo.RemoveAssociation(r.Context(), codPulv, codArea); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
