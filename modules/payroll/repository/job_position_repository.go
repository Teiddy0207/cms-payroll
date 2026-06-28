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


func (r *PayrollRepository) CreateJobPosition(ctx context.Context, pos *entity.JobPosition) (*entity.JobPosition, error) {
	if pos.ID == uuid.Nil {
		pos.ID = uuid.New()
	}
	query := `
		INSERT INTO job_descriptions (id, code, name, description, department_id, created_at, updated_at)
		VALUES (:id, :code, :name, :description, :department_id, NOW(), NOW())
		RETURNING id, code, name, description, department_id, created_at, updated_at
	`
	rows, err := r.DB.NamedQueryContext(ctx, query, pos)
	if err != nil {
		logger.Error("PayrollRepository:CreateJobPosition:Error %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdPos entity.JobPosition
		if err := rows.StructScan(&createdPos); err != nil {
			logger.Error("PayrollRepository:CreateJobPosition:Scan:Error %v", err)
			return nil, err
		}
		return &createdPos, nil
	}
	return nil, errors.New("failed to create job position")
}

func (r *PayrollRepository) GetJobPositions(ctx context.Context, qp params.QueryParams) ([]entity.JobPosition, int, error) {
	offset := (qp.PageNumber - 1) * qp.PageSize
	baseQuery := `FROM job_descriptions j`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if qp.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(j.name ILIKE $%d OR j.code ILIKE $%d)", argIndex, argIndex))
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
		logger.Error("PayrollRepository:GetJobPositions - Count", err)
		return nil, 0, err
	}

	dataQuery := `
		SELECT j.id, j.code, j.name, j.description, j.department_id, j.created_at, j.updated_at
	` + baseQuery + whereClause + ` ORDER BY j.name ASC, j.created_at DESC`

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, qp.PageSize, offset)

	var poses []entity.JobPosition
	err = r.DB.SelectContext(ctx, &poses, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.JobPosition{}, 0, nil
		}
		logger.Error("PayrollRepository:GetJobPositions - Select", err)
		return nil, 0, err
	}

	return poses, totalItems, nil
}

func (r *PayrollRepository) GetJobPositionByID(ctx context.Context, id uuid.UUID) (*entity.JobPosition, error) {
	var pos entity.JobPosition
	query := `SELECT id, code, name, description, department_id, created_at, updated_at FROM job_descriptions WHERE id = $1`
	err := r.DB.GetContext(ctx, &pos, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("PayrollRepository:GetJobPositionByID - Select", err)
		return nil, err
	}
	return &pos, nil
}

func (r *PayrollRepository) UpdateJobPosition(ctx context.Context, id uuid.UUID, pos *entity.JobPosition) error {
	query := `
		UPDATE job_descriptions
		SET name = :name, description = :description, department_id = :department_id, updated_at = NOW()
		WHERE id = :id
	`
	pos.ID = id
	result, err := r.DB.NamedExecContext(ctx, query, pos)
	if err != nil {
		logger.Error("PayrollRepository:UpdateJobPosition:Error %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("job position not found")
	}
	return nil
}

func (r *PayrollRepository) DeleteJobPosition(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM job_descriptions WHERE id = $1`
	_, err := r.DB.SQLx().ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("PayrollRepository:DeleteJobPosition:Error %v", err)
		return err
	}
	return nil
}
