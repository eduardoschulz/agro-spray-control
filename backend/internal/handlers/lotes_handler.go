package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/eduardoschulz/agro-spray-control/backend/internal/repositories"
	"github.com/gorilla/mux"
)

type LoteHandler struct {
	repo *repositories.LoteRepo
}

func NewLoteHandler(repo *repositories.LoteRepo) *LoteHandler {
	return &LoteHandler{repo: repo}
}

// CreateLote cria um novo lote
func (h *LoteHandler) CreateLote(w http.ResponseWriter, r *http.Request) {
	var lote models.Lote
	if err := json.NewDecoder(r.Body).Decode(&lote); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &lote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(lote)
}

// GetLote obtém um lote pelo código
func (h *LoteHandler) GetLote(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	lote, err := h.repo.GetByCod(r.Context(), cod)
	if err != nil {
		http.Error(w, "Lote não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lote)
}

// UpdateLote atualiza um lote existente
func (h *LoteHandler) UpdateLote(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	var lote models.Lote
	if err := json.NewDecoder(r.Body).Decode(&lote); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garante que o código na URL corresponde ao corpo da requisição
	lote.Cod = cod

	if err := h.repo.Update(r.Context(), &lote); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lote)
}

// DeleteLote remove um lote
func (h *LoteHandler) DeleteLote(w http.ResponseWriter, r *http.Request) {
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

// ListLotes lista todos os lotes
func (h *LoteHandler) ListLotes(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	lotes, err := h.repo.List(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lotes)
}
