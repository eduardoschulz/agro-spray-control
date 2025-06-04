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

func (r *UsuarioRepo) Update(ctx context.Context, usuario *models.Usuario) error {
	query := `
		UPDATE usuarios
		SET nome = :nome,
			email = :email,
			password_hash = :password_hash,
			nivel_permissao = :nivel_permissao
		WHERE cpf = :cpf
	`

	_, err := r.db.NamedExecContext(ctx, query, usuario)
	if err != nil {
		log.Printf("Erro ao atualizar usuário: %v", err)
		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}
	return nil
}

func (r *UsuarioRepo) Delete(ctx context.Context, cpf string) error {
	query := `
		DELETE FROM usuarios
		WHERE cpf = $1
	`

	_, err := r.db.ExecContext(ctx, query, cpf)
	if err != nil {
		log.Printf("Erro ao deletar usuário: %v", err)
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}
	return nil
}

func (r *UsuarioRepo) List(ctx context.Context, page, limit int) ([]models.Usuario, error) {
	offset := (page - 1) * limit
	query := `
		SELECT cpf, nome, email, password_hash, nivel_permissao, criado_em
		FROM usuarios
		ORDER BY criado_em DESC
		LIMIT $1 OFFSET $2
	`

	var usuarios []models.Usuario
	err := r.db.SelectContext(ctx, &usuarios, query, limit, offset)
	if err != nil {
		log.Printf("Erro ao listar usuários: %v", err)
		return nil, fmt.Errorf("erro ao listar usuários: %w", err)
	}

	return usuarios, nil
}


