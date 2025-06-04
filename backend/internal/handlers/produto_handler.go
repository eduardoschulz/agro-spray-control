package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/eduardoschulz/agro-spray-control/backend/internal/repositories"
	"github.com/gorilla/mux"
)

type ProdutoHandler struct {
	repo *repositories.ProdutoRepo
}

func NewProdutoHandler(repo *repositories.ProdutoRepo) *ProdutoHandler {
	return &ProdutoHandler{repo: repo}
}

// CreateProduto cria um novo produto
func (h *ProdutoHandler) CreateProduto(w http.ResponseWriter, r *http.Request) {
	var produto models.Produto
	if err := json.NewDecoder(r.Body).Decode(&produto); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &produto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(produto)
}

// GetProduto obtém um produto pelo código
func (h *ProdutoHandler) GetProduto(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	produto, err := h.repo.GetByCod(r.Context(), cod)
	if err != nil {
		http.Error(w, "Produto não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produto)
}

// UpdateProduto atualiza um produto existente
func (h *ProdutoHandler) UpdateProduto(w http.ResponseWriter, r *http.Request) {
	cod := mux.Vars(r)["cod"]
	if cod == "" {
		http.Error(w, "Código é obrigatório", http.StatusBadRequest)
		return
	}

	var produto models.Produto
	if err := json.NewDecoder(r.Body).Decode(&produto); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garante que o código na URL corresponde ao corpo da requisição
	produto.Cod = cod

	if err := h.repo.Update(r.Context(), &produto); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produto)
}

// DeleteProduto remove um produto
func (h *ProdutoHandler) DeleteProduto(w http.ResponseWriter, r *http.Request) {
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

// ListProdutos lista todos os produtos
func (h *ProdutoHandler) ListProdutos(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	produtos, err := h.repo.List(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(produtos)
}
