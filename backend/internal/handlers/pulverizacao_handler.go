package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/eduardoschulz/agro-spray-control/backend/internal/repositories"
	"github.com/gorilla/mux"
)

type PulverizacaoHandler struct {
	repo *repositories.PulverizacaoRepo
}

func NewPulverizacaoHandler(repo *repositories.PulverizacaoRepo) *PulverizacaoHandler {
	return &PulverizacaoHandler{repo: repo}
}

// CreatePulverizacao cria uma nova pulverização
func (h *PulverizacaoHandler) CreatePulverizacao(w http.ResponseWriter, r *http.Request) {
	var pulverizacao models.Pulverizacao
	if err := json.NewDecoder(r.Body).Decode(&pulverizacao); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &pulverizacao); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(pulverizacao)
}

// GetPulverizacao obtém uma pulverização pelo código
func (h *PulverizacaoHandler) GetPulverizacao(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	pulverizacao, err := h.repo.GetByCod(r.Context(), cod)
	if err != nil {
		http.Error(w, "Pulverização não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pulverizacao)
}

// UpdatePulverizacao atualiza uma pulverização existente
func (h *PulverizacaoHandler) UpdatePulverizacao(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	var pulverizacao models.Pulverizacao
	if err := json.NewDecoder(r.Body).Decode(&pulverizacao); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garante que o código na URL corresponde ao corpo da requisição
	pulverizacao.Cod = cod

	if err := h.repo.Update(r.Context(), &pulverizacao); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pulverizacao)
}

// DeletePulverizacao remove uma pulverização
func (h *PulverizacaoHandler) DeletePulverizacao(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), cod); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListPulverizacoes lista todas as pulverizações
func (h *PulverizacaoHandler) ListPulverizacoes(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	pulverizacoes, err := h.repo.List(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pulverizacoes)
}
