package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/eduardoschulz/agro-spray-control/backend/internal/repositories"
	"github.com/gorilla/mux"
)

type AreaHandler struct {
	repo *repositories.AreaRepo
}

func NewAreaHandler(repo *repositories.AreaRepo) *AreaHandler {
	return &AreaHandler{repo: repo}
}

// CreateArea cria uma nova área
func (h *AreaHandler) CreateArea(w http.ResponseWriter, r *http.Request) {
	var area models.Area
	if err := json.NewDecoder(r.Body).Decode(&area); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &area); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(area)
}

// GetArea obtém uma área pelo código
func (h *AreaHandler) GetArea(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	area, err := h.repo.GetByCod(r.Context(), cod)
	if err != nil {
		http.Error(w, "Área não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(area)
}

// UpdateArea atualiza uma área existente
func (h *AreaHandler) UpdateArea(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	var area models.Area
	if err := json.NewDecoder(r.Body).Decode(&area); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garante que o código na URL corresponde ao corpo da requisição
	area.Cod = cod

	if err := h.repo.Update(r.Context(), &area); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(area)
}

// DeleteArea remove uma área
func (h *AreaHandler) DeleteArea(w http.ResponseWriter, r *http.Request) {
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

// ListAreas lista todas as áreas
func (h *AreaHandler) ListAreas(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	areas, err := h.repo.List(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(areas)
}
