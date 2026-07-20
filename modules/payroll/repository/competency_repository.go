package repository

import (
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/payroll/entity"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// ==========================================
// Competency Dictionaries Repository Implementation
// ==========================================

func (r *PayrollRepository) CreateCompetency(ctx context.Context, c *entity.Competency) (*entity.Competency, error) {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	query := `
		INSERT INTO competency_dictionaries (id, code, name, point_value, description, type_id, created_at, updated_at)
		VALUES (:id, :code, :name, :point_value, :description, :type_id, NOW(), NOW())
		RETURNING id, code, name, point_value, description, type_id, created_at, updated_at
	`
	rows, err := r.DB.NamedQueryContext(ctx, query, c)
	if err != nil {
		logger.Error("PayrollRepository:CreateCompetency:Error %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var created entity.Competency
		if err := rows.StructScan(&created); err != nil {
			logger.Error("PayrollRepository:CreateCompetency:Scan:Error %v", err)
			return nil, err
		}
		return &created, nil
	}
	return nil, errors.New("failed to create competency")
}

func (r *PayrollRepository) GetCompetencies(ctx context.Context, qp params.QueryParams) ([]entity.Competency, int, error) {
	offset := (qp.PageNumber - 1) * qp.PageSize
	baseQuery := `FROM competency_dictionaries c`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if qp.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(c.name ILIKE $%d OR c.code ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+qp.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = " WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + baseQuery + whereClause
	var totalItems int
	err := r.DB.GetContext(ctx, &totalItems, countQuery, args...)
	if err != nil {
		logger.Error("PayrollRepository:GetCompetencies - Count", err)
		return nil, 0, err
	}

	dataQuery := `
		SELECT c.id, c.code, c.name, c.point_value, c.description, c.type_id, c.created_at, c.updated_at
	` + baseQuery + whereClause + ` ORDER BY c.code ASC, c.created_at DESC`

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, qp.PageSize, offset)

	var list []entity.Competency
	err = r.DB.SelectContext(ctx, &list, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.Competency{}, 0, nil
		}
		logger.Error("PayrollRepository:GetCompetencies - Select", err)
		return nil, 0, err
	}

	return list, totalItems, nil
}

func (r *PayrollRepository) GetCompetencyByID(ctx context.Context, id uuid.UUID) (*entity.Competency, error) {
	var c entity.Competency
	query := `SELECT * FROM competency_dictionaries WHERE id = $1`
	err := r.DB.GetContext(ctx, &c, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("PayrollRepository:GetCompetencyByID - Select", err)
		return nil, err
	}
	return &c, nil
}

func (r *PayrollRepository) UpdateCompetency(ctx context.Context, id uuid.UUID, c *entity.Competency) error {
	query := `
		UPDATE competency_dictionaries
		SET name = :name, point_value = :point_value, description = :description, type_id = :type_id, updated_at = NOW()
		WHERE id = :id
	`
	c.ID = id
	result, err := r.DB.NamedExecContext(ctx, query, c)
	if err != nil {
		logger.Error("PayrollRepository:UpdateCompetency:Error %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("competency not found")
	}
	return nil
}

func (r *PayrollRepository) DeleteCompetency(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM competency_dictionaries WHERE id = $1`
	_, err := r.DB.SQLx().ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("PayrollRepository:DeleteCompetency:Error %v", err)
		return err
	}
	return nil
}
