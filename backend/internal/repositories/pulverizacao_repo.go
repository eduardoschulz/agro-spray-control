package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type PulverizacaoRepo struct {
	db *sqlx.DB
}

func NewPulverizacaoRepo(db *sqlx.DB) *PulverizacaoRepo {
	return &PulverizacaoRepo{db: db}
}

func (r *PulverizacaoRepo) Create(ctx context.Context, pulverizacao *models.Pulverizacao) error {
	query := `
		INSERT INTO pulverizacao (dtaplicacao, cultura, cod_lote, cpf_responsavel)
		VALUES (:dt_aplicacao, :cultura, :cod_lote, :cpf_responsavel)
		RETURNING cod
	`
	rows, err := r.db.NamedQueryContext(ctx, query, pulverizacao)
	if err != nil {
		return fmt.Errorf("erro ao criar pulverização: %w", err)
	}
	
	if rows.Next() {
		return rows.Scan(&pulverizacao.Cod)
	}
	return rows.Close()
}

func (r *PulverizacaoRepo) GetByCod(ctx context.Context, cod string) (*models.Pulverizacao, error) {
	query := `
		SELECT cod, dtaplicacao as dt_aplicacao, cultura, cod_lote, cpf_responsavel
		FROM pulverizacao
		WHERE cod = $1
	`
	var pulverizacao models.Pulverizacao
	err := r.db.GetContext(ctx, &pulverizacao, query, cod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("pulverização não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar pulverização: %w", err)
	}
	return &pulverizacao, nil
}

func (r *PulverizacaoRepo) Update(ctx context.Context, pulverizacao *models.Pulverizacao) error {
	query := `
		UPDATE pulverizacao SET
			dtaplicacao = :dt_aplicacao,
			cultura = :cultura,
			cod_lote = :cod_lote,
			cpf_responsavel = :cpf_responsavel
		WHERE cod = :cod
	`
	result, err := r.db.NamedExecContext(ctx, query, pulverizacao)
	if err != nil {
		return fmt.Errorf("erro ao atualizar pulverização: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhuma pulverização foi atualizada")
	}
	return nil
}

func (r *PulverizacaoRepo) Delete(ctx context.Context, cod string) error {
	query := `DELETE FROM pulverizacao WHERE cod = $1`
	result, err := r.db.ExecContext(ctx, query, cod)
	if err != nil {
		return fmt.Errorf("erro ao deletar pulverização: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhuma pulverização foi deletada")
	}
	return nil
}

func (r *PulverizacaoRepo) List(ctx context.Context, page, limit int) ([]models.Pulverizacao, error) {
	offset := (page - 1) * limit
	query := `
		SELECT cod, dtaplicacao as dt_aplicacao, cultura, cod_lote, cpf_responsavel
		FROM pulverizacao
		ORDER BY dtaplicacao DESC
		LIMIT $1 OFFSET $2
	`
	var pulverizacoes []models.Pulverizacao
	err := r.db.SelectContext(ctx, &pulverizacoes, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pulverizações: %w", err)
	}
	return pulverizacoes, nil
}

func (r *PulverizacaoRepo) ListByResponsavel(ctx context.Context, cpfResponsavel string) ([]models.Pulverizacao, error) {
	query := `
		SELECT cod, dtaplicacao as dt_aplicacao, cultura, cod_lote, cpf_responsavel
		FROM pulverizacao
		WHERE cpf_responsavel = $1
		ORDER BY dtaplicacao DESC
	`
	var pulverizacoes []models.Pulverizacao
	err := r.db.SelectContext(ctx, &pulverizacoes, query, cpfResponsavel)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pulverizações por responsável: %w", err)
	}
	return pulverizacoes, nil
}

func (r *PulverizacaoRepo) ListByLote(ctx context.Context, codLote string) ([]models.Pulverizacao, error) {
	query := `
		SELECT cod, dtaplicacao as dt_aplicacao, cultura, cod_lote, cpf_responsavel
		FROM pulverizacao
		WHERE cod_lote = $1
		ORDER BY dtaplicacao DESC
	`
	var pulverizacoes []models.Pulverizacao
	err := r.db.SelectContext(ctx, &pulverizacoes, query, codLote)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar pulverizações por lote: %w", err)
	}
	return pulverizacoes, nil
}
