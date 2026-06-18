package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/m0hit2409/transactions-service/internal/domain"
)

type operationTypeRepo struct {
	db *sqlx.DB
}

// NewOperationTypeRepository returns a SQLite-backed OperationTypeRepository.
func NewOperationTypeRepository(db *sqlx.DB) OperationTypeRepository {
	return &operationTypeRepo{db: db}
}

func (r *operationTypeRepo) FindByID(ctx context.Context, id int) (domain.OperationType, error) {
	const q = `
		SELECT operation_type_id, description, sign, is_active, created_at
		FROM operation_types
		WHERE operation_type_id = ?`

	var ot domain.OperationType
	if err := r.db.GetContext(ctx, &ot, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.OperationType{}, domain.ErrOperationTypeNotFound
		}
		return domain.OperationType{}, fmt.Errorf("select operation type: %w", err)
	}
	return ot, nil
}
