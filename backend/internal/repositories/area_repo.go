package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type AreaRepo struct {
	db *sqlx.DB
}

func NewAreaRepo(db *sqlx.DB) *AreaRepo {
	return &AreaRepo{db: db}
}

func (r *AreaRepo) Create(ctx context.Context, area *models.Area) error {
	query := `
		INSERT INTO area (tamanho, fazenda_cod)
		VALUES (:tamanho, :cod_fazenda)
		RETURNING cod
	`
	rows, err := r.db.NamedQueryContext(ctx, query, area)
	if err != nil {
		return fmt.Errorf("erro ao criar área: %w", err)
	}
	
	if rows.Next() {
		return rows.Scan(&area.Cod)
	}
	return rows.Close()
}

func (r *AreaRepo) GetByCod(ctx context.Context, cod string) (*models.Area, error) {
	query := `
		SELECT cod, tamanho, fazenda_cod as cod_fazenda
		FROM area
		WHERE cod = $1
	`
	var area models.Area
	err := r.db.GetContext(ctx, &area, query, cod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("área não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar área: %w", err)
	}
	return &area, nil
}

func (r *AreaRepo) Update(ctx context.Context, area *models.Area) error {
	query := `
		UPDATE area SET
			tamanho = :tamanho,
			fazenda_cod = :cod_fazenda
		WHERE cod = :cod
	`
	result, err := r.db.NamedExecContext(ctx, query, area)
	if err != nil {
		return fmt.Errorf("erro ao atualizar área: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhuma área foi atualizada")
	}
	return nil
}

func (r *AreaRepo) Delete(ctx context.Context, cod string) error {
	query := `DELETE FROM area WHERE cod = $1`
	result, err := r.db.ExecContext(ctx, query, cod)
	if err != nil {
		return fmt.Errorf("erro ao deletar área: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhuma área foi deletada")
	}
	return nil
}

func (r *AreaRepo) List(ctx context.Context, page, limit int) ([]models.Area, error) {
	offset := (page - 1) * limit
	query := `
		SELECT cod, tamanho, fazenda_cod as cod_fazenda
		FROM area
		ORDER BY cod
		LIMIT $1 OFFSET $2
	`
	var areas []models.Area
	err := r.db.SelectContext(ctx, &areas, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar áreas: %w", err)
	}
	return areas, nil
}

func (r *AreaRepo) ListByFazenda(ctx context.Context, codFazenda string) ([]models.Area, error) {
	query := `
		SELECT cod, tamanho, fazenda_cod as cod_fazenda
		FROM area
		WHERE fazenda_cod = $1
		ORDER BY cod
	`
	var areas []models.Area
	err := r.db.SelectContext(ctx, &areas, query, codFazenda)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar áreas por fazenda: %w", err)
	}
	return areas, nil
}
