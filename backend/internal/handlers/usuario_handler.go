package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
    "time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/eduardoschulz/agro-spray-control/backend/internal/repositories"
	"github.com/gorilla/mux"
)

var jwtKey = []byte("senhaforte123")

type UsuarioHandler struct {
	repo *repositories.UsuarioRepo
}

func NewUsuarioHandler(repo *repositories.UsuarioRepo) *UsuarioHandler {
	return &UsuarioHandler{repo: repo}
}

// CreateUsuario cria um novo usuário
func (h *UsuarioHandler) CreateUsuario(w http.ResponseWriter, r *http.Request) {
	var usuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(usuario)
}

// GetUsuario obtém um usuário pelo CPF
func (h *UsuarioHandler) GetUsuario(w http.ResponseWriter, r *http.Request) {
	cpf := mux.Vars(r)["cpf"]
	if cpf == "" {
		http.Error(w, "CPF é obrigatório", http.StatusBadRequest)
		return
	}

	usuario, err := h.repo.GetByCPF(r.Context(), cpf)
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// UpdateUsuario atualiza um usuário existente
func (h *UsuarioHandler) UpdateUsuario(w http.ResponseWriter, r *http.Request) {
	cpf := mux.Vars(r)["cpf"]
	if cpf == "" {
		http.Error(w, "CPF é obrigatório", http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := json.NewDecoder(r.Body).Decode(&usuario); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	// Garante que o CPF na URL corresponde ao corpo da requisição
	usuario.CPF = cpf

	if err := h.repo.Update(r.Context(), &usuario); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

// DeleteUsuario remove um usuário
func (h *UsuarioHandler) DeleteUsuario(w http.ResponseWriter, r *http.Request) {
	cpf := mux.Vars(r)["cpf"]
	if cpf == "" {
		http.Error(w, "CPF é obrigatório", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), cpf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListUsuarios lista todos os usuários
func (h *UsuarioHandler) ListUsuarios(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	page, _ := strconv.Atoi(queryValues.Get("page"))
	limit, _ := strconv.Atoi(queryValues.Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	usuarios, err := h.repo.List(r.Context(), page, limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}

func (h *UsuarioHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CPF   string `json:"cpf"`
		Senha string `json:"senha"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Dados inválidos", http.StatusBadRequest)
		return
	}

	usuario, err := h.repo.GetByCPF(r.Context(), input.CPF)
	if err != nil {
		http.Error(w, "Usuário não encontrado", http.StatusUnauthorized)
		return
	}

	if input.Senha != string(usuario.PasswordHash) {
		http.Error(w, "Senha incorreta", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"cpf": usuario.CPF,
		"exp": time.Now().Add(2 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Erro ao gerar token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}
