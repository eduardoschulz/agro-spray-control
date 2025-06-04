package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/eduardoschulz/agro-spray-control/backend/internal/models"
	"github.com/jmoiron/sqlx"
)

type FazendaRepo struct {
	db *sqlx.DB
}

func NewFazendaRepo(db *sqlx.DB) *FazendaRepo {
	return &FazendaRepo{db: db}
}

func (r *FazendaRepo) Create(ctx context.Context, fazenda *models.Fazenda) error {
	query := `
		INSERT INTO fazenda (localizacao, cpf_proprietario)
		VALUES (:localizacao, :cpf_proprietario)
		RETURNING cod
	`
	rows, err := r.db.NamedQueryContext(ctx, query, fazenda)
	if err != nil {
		return fmt.Errorf("erro ao criar fazenda: %w", err)
	}
	
	if rows.Next() {
		return rows.Scan(&fazenda.Cod)
	}
	return rows.Close()
}

func (r *FazendaRepo) GetByCod(ctx context.Context, cod string) (*models.Fazenda, error) {
	query := `
		SELECT cod, localizacao, cpf_proprietario
		FROM fazenda
		WHERE cod = $1
	`
	var fazenda models.Fazenda
	err := r.db.GetContext(ctx, &fazenda, query, cod)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("fazenda não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar fazenda: %w", err)
	}
	return &fazenda, nil
}

func (r *FazendaRepo) Update(ctx context.Context, fazenda *models.Fazenda) error {
	query := `
		UPDATE fazenda SET
			localizacao = :localizacao,
			cpf_proprietario = :cpf_proprietario
		WHERE cod = :cod
	`
	result, err := r.db.NamedExecContext(ctx, query, fazenda)
	if err != nil {
		return fmt.Errorf("erro ao atualizar fazenda: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhuma fazenda foi atualizada")
	}
	return nil
}

func (r *FazendaRepo) Delete(ctx context.Context, cod string) error {
	query := `DELETE FROM fazenda WHERE cod = $1`
	result, err := r.db.ExecContext(ctx, query, cod)
	if err != nil {
		return fmt.Errorf("erro ao deletar fazenda: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("nenhuma fazenda foi deletada")
	}
	return nil
}

func (r *FazendaRepo) List(ctx context.Context, page, limit int) ([]models.Fazenda, error) {
	offset := (page - 1) * limit
	query := `
		SELECT cod, localizacao, cpf_proprietario
		FROM fazenda
		ORDER BY cod
		LIMIT $1 OFFSET $2
	`
	var fazendas []models.Fazenda
	err := r.db.SelectContext(ctx, &fazendas, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar fazendas: %w", err)
	}
	return fazendas, nil
}

func (r *FazendaRepo) ListByProprietario(ctx context.Context, cpfProprietario string) ([]models.Fazenda, error) {
	query := `
		SELECT cod, localizacao, cpf_proprietario
		FROM fazenda
		WHERE cpf_proprietario = $1
		ORDER BY cod
	`
	var fazendas []models.Fazenda
	err := r.db.SelectContext(ctx, &fazendas, query, cpfProprietario)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar fazendas por proprietário: %w", err)
	}
	return fazendas, nil
}
