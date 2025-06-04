package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type PulverizacaoAreaRepo struct {
	db *sqlx.DB
}

func NewPulverizacaoAreaRepo(db *sqlx.DB) *PulverizacaoAreaRepo {
	return &PulverizacaoAreaRepo{db: db}
}

func (r *PulverizacaoAreaRepo) Associate(ctx context.Context, pa *models.PulverizacaoArea) error {
	query := `
		INSERT INTO pulverizacao_areas (cod_pulv, cod_area)
		VALUES (:cod_pulv, :cod_area)
	`
	_, err := r.db.NamedExecContext(ctx, query, pa)
	if err != nil {
		return fmt.Errorf("erro ao associar área à pulverização: %w", err)
	}
	return nil
}

func (r *PulverizacaoAreaRepo) GetAreasByPulverizacao(ctx context.Context, codPulv string) ([]models.Area, error) {
	query := `
		SELECT a.cod, a.tamanho, a.fazenda_cod as cod_fazenda
		FROM area a
		JOIN pulverizacao_areas pa ON a.cod = pa.cod_area
		WHERE pa.cod_pulv = $1
	`
	var areas []models.Area
	err := r.db.SelectContext(ctx, &areas, query, codPulv)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar áreas da pulverização: %w", err)
	}
	return areas, nil
}

func (r *PulverizacaoAreaRepo) GetPulverizacoesByArea(ctx context.Context, codArea string) ([]models.Pulverizacao, error) {
	query := `
		SELECT p.cod, p.dtaplicacao as dt_aplicacao, p.cultura, p.cod_lote, p.cpf_responsavel
		FROM pulverizacao p
		JOIN pulverizacao_areas pa ON p.cod = pa.cod_pulv
		WHERE pa.cod_area = $1
		ORDER BY p.dtaplicacao DESC
	`
	var pulverizacoes []models.Pulverizacao
	err := r.db.SelectContext(ctx, &pulverizacoes, query, codArea)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar pulverizações da área: %w", err)
	}
	return pulverizacoes, nil
}

func (r *PulverizacaoAreaRepo) RemoveAssociation(ctx context.Context, codPulv, codArea string) error {
	query := `
		DELETE FROM pulverizacao_areas
		WHERE cod_pulv = $1 AND cod_area = $2
	`
	result, err := r.db.ExecContext(ctx, query, codPulv, codArea)
	if err != nil {
		return fmt.Errorf("erro ao remover associação: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhuma associação foi removida")
	}
	return nil
}
