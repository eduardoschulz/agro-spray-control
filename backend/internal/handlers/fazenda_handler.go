package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/eduardoschulz/agro-spray-control/backend/internal/repositories"
	"github.com/gorilla/mux"
)

type FazendaHandler struct {
	repo *repositories.FazendaRepo
}

func NewFazendaHandler(repo *repositories.FazendaRepo) *FazendaHandler {
	return &FazendaHandler{repo: repo}
}

// CreateFazenda cria uma nova fazenda
func (h *FazendaHandler) CreateFazenda(w http.ResponseWriter, r *http.Request) {
	var fazenda models.Fazenda
	if err := json.NewDecoder(r.Body).Decode(&fazenda); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &fazenda); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(fazenda)
}

// GetFazenda obtém uma fazenda pelo código
func (h *FazendaHandler) GetFazenda(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	fazenda, err := h.repo.GetByCod(r.Context(), cod)
	if err != nil {
		http.Error(w, "Fazenda não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fazenda)
}

// UpdateFazenda atualiza uma fazenda existente
func (h *FazendaHandler) UpdateFazenda(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	var fazenda models.Fazenda
	if err := json.NewDecoder(r.Body).Decode(&fazenda); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garante que o código na URL corresponde ao corpo da requisição
	fazenda.Cod = cod

	if err := h.repo.Update(r.Context(), &fazenda); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fazenda)
}

// DeleteFazenda remove uma fazenda
func (h *FazendaHandler) DeleteFazenda(w http.ResponseWriter, r *http.Request) {
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

// ListFazendas lista todas as fazendas
func (h *FazendaHandler) ListFazendas(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	fazendas, err := h.repo.List(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fazendas)
}
