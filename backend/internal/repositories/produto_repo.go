package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type ProdutoRepo struct {
	db *sqlx.DB
}

func NewProdutoRepo(db *sqlx.DB) *ProdutoRepo {
	return &ProdutoRepo{db: db}
}

func (r *ProdutoRepo) Create(ctx context.Context, produto *models.Produto) error {
	query := `
		INSERT INTO produtos (descricao, fabricante, compquimica)
		VALUES (:descricao, :fabricante, :comp_quimica)
		RETURNING cod
	`
	rows, err := r.db.NamedQueryContext(ctx, query, produto)
	if err != nil {
		return fmt.Errorf("erro ao criar produto: %w", err)
	}
	
	if rows.Next() {
		return rows.Scan(&produto.Cod)
	}
	return rows.Close()
}

func (r *ProdutoRepo) GetByCod(ctx context.Context, cod string) (*models.Produto, error) {
	query := `
		SELECT cod, descricao, fabricante, compquimica as comp_quimica
		FROM produtos
		WHERE cod = $1
	`
	var produto models.Produto
	err := r.db.GetContext(ctx, &produto, query, cod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("produto não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}
	return &produto, nil
}

func (r *ProdutoRepo) Update(ctx context.Context, produto *models.Produto) error {
	query := `
		UPDATE produtos SET
			descricao = :descricao,
			fabricante = :fabricante,
			compquimica = :comp_quimica
		WHERE cod = :cod
	`
	result, err := r.db.NamedExecContext(ctx, query, produto)
	if err != nil {
		return fmt.Errorf("erro ao atualizar produto: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhum produto foi atualizado")
	}
	return nil
}

func (r *ProdutoRepo) Delete(ctx context.Context, cod string) error {
	query := `DELETE FROM produtos WHERE cod = $1`
	result, err := r.db.ExecContext(ctx, query, cod)
	if err != nil {
		return fmt.Errorf("erro ao deletar produto: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhum produto foi deletado")
	}
	return nil
}

func (r *ProdutoRepo) List(ctx context.Context, page, limit int) ([]models.Produto, error) {
	offset := (page - 1) * limit
	query := `
		SELECT cod, descricao, fabricante, compquimica as comp_quimica
		FROM produtos
		ORDER BY descricao
		LIMIT $1 OFFSET $2
	`
	var produtos []models.Produto
	err := r.db.SelectContext(ctx, &produtos, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar produtos: %w", err)
	}
	return produtos, nil
}
