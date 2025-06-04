package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type LoteRepo struct {
	db *sqlx.DB
}

func NewLoteRepo(db *sqlx.DB) *LoteRepo {
	return &LoteRepo{db: db}
}

func (r *LoteRepo) Create(ctx context.Context, lote *models.Lote) error {
	query := `
		INSERT INTO lotes (dtvalidade, cod_produto, quantidade)
		VALUES (:dt_validade, :cod_produto, :quantidade)
		RETURNING cod
	`
	rows, err := r.db.NamedQueryContext(ctx, query, lote)
	if err != nil {
		return fmt.Errorf("erro ao criar lote: %w", err)
	}
	
	if rows.Next() {
		return rows.Scan(&lote.Cod)
	}
	return rows.Close()
}

func (r *LoteRepo) GetByCod(ctx context.Context, cod string) (*models.Lote, error) {
	query := `
		SELECT cod, dtvalidade as dt_validade, quantidade, cod_produto
		FROM lotes
		WHERE cod = $1
	`
	var lote models.Lote
	err := r.db.GetContext(ctx, &lote, query, cod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("lote não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar lote: %w", err)
	}
	return &lote, nil
}

func (r *LoteRepo) Update(ctx context.Context, lote *models.Lote) error {
	query := `
		UPDATE lotes SET
			dtvalidade = :dt_validade,
			cod_produto = :cod_produto,
			quantidade = :quantidade
		WHERE cod = :cod
	`
	result, err := r.db.NamedExecContext(ctx, query, lote)
	if err != nil {
		return fmt.Errorf("erro ao atualizar lote: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhum lote foi atualizado")
	}
	return nil
}

func (r *LoteRepo) Delete(ctx context.Context, cod string) error {
	query := `DELETE FROM lotes WHERE cod = $1`
	result, err := r.db.ExecContext(ctx, query, cod)
	if err != nil {
		return fmt.Errorf("erro ao deletar lote: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhum lote foi deletado")
	}
	return nil
}

func (r *LoteRepo) List(ctx context.Context, page, limit int) ([]models.Lote, error) {
	offset := (page - 1) * limit
	query := `
		SELECT cod, dtvalidade as dt_validade, quantidade, cod_produto
		FROM lotes
		ORDER BY dtvalidade
		LIMIT $1 OFFSET $2
	`
	var lotes []models.Lote
	err := r.db.SelectContext(ctx, &lotes, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar lotes: %w", err)
	}
	return lotes, nil
}

func (r *LoteRepo) ListByProduto(ctx context.Context, codProduto string) ([]models.Lote, error) {
	query := `
		SELECT cod, dtvalidade as dt_validade, quantidade, cod_produto
		FROM lotes
		WHERE cod_produto = $1
		ORDER BY dtvalidade
	`
	var lotes []models.Lote
	err := r.db.SelectContext(ctx, &lotes, query, codProduto)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar lotes por produto: %w", err)
	}
	return lotes, nil
}

func (r *LoteRepo) ListVencidos(ctx context.Context) ([]models.Lote, error) {
	query := `
		SELECT cod, dtvalidade as dt_validade, quantidade, cod_produto
		FROM lotes
		WHERE dtvalidade < $1
		ORDER BY dtvalidade
	`
	var lotes []models.Lote
	err := r.db.SelectContext(ctx, &lotes, query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("erro ao listar lotes vencidos: %w", err)
	}
	return lotes, nil
}
