package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UsuarioRepo struct {
    db *sqlx.DB
}

func NewUsuarioRepo(db *sqlx.DB) *UsuarioRepo {
    return &UsuarioRepo{db: db}
}

func (r *UsuarioRepo) Create(ctx context.Context, usuario *models.Usuario) error {
    query := `
        INSERT INTO usuarios (cpf, nome, email, password_hash, nivel_permissao, criado_em)
        VALUES (:cpf, :nome, :email, :password_hash, :nivel_permissao, :criado_em)
        RETURNING cpf
    `
    
    // Usando NamedQuery do SQLx
    rows, err := r.db.NamedQueryContext(ctx, query, usuario)
    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok {
            switch pqErr.Code.Name() {
            case "unique_violation":
                if pqErr.Constraint == "usuarios_pkey" {
                    return errors.New("CPF já cadastrado")
                } else if pqErr.Constraint == "usuarios_email_key" {
                    return errors.New("e-mail já cadastrado")
                }
            }
        }
        log.Printf("Erro ao criar usuário: %v", err)
        return fmt.Errorf("erro ao criar usuário: %w", err)
    }
    
    if rows.Next() {
        rows.Scan(&usuario.CPF)
    }
    return rows.Close()
}

func (r *UsuarioRepo) GetByCPF(ctx context.Context, cpf string) (*models.Usuario, error) {
    query := `
        SELECT cpf, nome, email, password_hash, nivel_permissao, criado_em
        FROM usuarios
        WHERE cpf = $1
    `
    var usuario models.Usuario
    err := r.db.GetContext(ctx, &usuario, query, cpf)
    if err != nil {
        return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
    }
    return &usuario, nil
}

// Adicionar métodos Update, Delete, List, etc.
