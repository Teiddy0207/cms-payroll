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
// Job Standards Repository Implementation
// ==========================================

func (r *PayrollRepository) CreateJobStandard(ctx context.Context, js *entity.JobStandard) (*entity.JobStandard, error) {
	if js.ID == uuid.Nil {
		js.ID = uuid.New()
	}
	query := `
		INSERT INTO job_standards (id, standard_code, name, allowance_value, description, created_at, updated_at)
		VALUES (:id, :standard_code, :name, :allowance_value, :description, NOW(), NOW())
		RETURNING id, standard_code, name, allowance_value, description, created_at, updated_at
	`
	rows, err := r.DB.NamedQueryContext(ctx, query, js)
	if err != nil {
		logger.Error("PayrollRepository:CreateJobStandard:Error %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var created entity.JobStandard
		if err := rows.StructScan(&created); err != nil {
			logger.Error("PayrollRepository:CreateJobStandard:Scan:Error %v", err)
			return nil, err
		}
		return &created, nil
	}
	return nil, errors.New("failed to create job standard")
}

func (r *PayrollRepository) GetJobStandards(ctx context.Context, qp params.QueryParams) ([]entity.JobStandard, int, error) {
	offset := (qp.PageNumber - 1) * qp.PageSize
	baseQuery := `FROM job_standards j`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if qp.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(j.name ILIKE $%d OR j.standard_code ILIKE $%d)", argIndex, argIndex))
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
		logger.Error("PayrollRepository:GetJobStandards - Count", err)
		return nil, 0, err
	}

	dataQuery := `
		SELECT j.id, j.standard_code, j.name, j.allowance_value, j.description, j.created_at, j.updated_at
	` + baseQuery + whereClause + ` ORDER BY j.name ASC, j.created_at DESC`

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, qp.PageSize, offset)

	var standards []entity.JobStandard
	err = r.DB.SelectContext(ctx, &standards, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.JobStandard{}, 0, nil
		}
		logger.Error("PayrollRepository:GetJobStandards - Select", err)
		return nil, 0, err
	}

	return standards, totalItems, nil
}

func (r *PayrollRepository) GetJobStandardByID(ctx context.Context, id uuid.UUID) (*entity.JobStandard, error) {
	var js entity.JobStandard
	query := `SELECT * FROM job_standards WHERE id = $1`
	err := r.DB.GetContext(ctx, &js, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("PayrollRepository:GetJobStandardByID - Select", err)
		return nil, err
	}
	return &js, nil
}

func (r *PayrollRepository) UpdateJobStandard(ctx context.Context, id uuid.UUID, js *entity.JobStandard) error {
	query := `
		UPDATE job_standards
		SET name = :name, allowance_value = :allowance_value, description = :description, updated_at = NOW()
		WHERE id = :id
	`
	js.ID = id
	result, err := r.DB.NamedExecContext(ctx, query, js)
	if err != nil {
		logger.Error("PayrollRepository:UpdateJobStandard:Error %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("job standard not found")
	}
	return nil
}

func (r *PayrollRepository) DeleteJobStandard(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM job_standards WHERE id = $1`
	_, err := r.DB.SQLx().ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("PayrollRepository:DeleteJobStandard:Error %v", err)
		return err
	}
	return nil
}
