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
		INSERT INTO job_descriptions (
			id, code, name, description, department_id,
			e_score, c_score, r_score, we_weight, wc_weight, wr_weight,
			salary_spread, job_score, midpoint, min_salary, max_salary,
			is_benchmark, market_salary, search_keyword,
			created_at, updated_at
		)
		VALUES (
			:id, :code, :name, :description, :department_id,
			:e_score, :c_score, :r_score, :we_weight, :wc_weight, :wr_weight,
			:salary_spread, :job_score, :midpoint, :min_salary, :max_salary,
			:is_benchmark, :market_salary, :search_keyword,
			NOW(), NOW()
		)
		RETURNING 
			id, code, name, description, department_id,
			e_score, c_score, r_score, we_weight, wc_weight, wr_weight,
			salary_spread, job_score, midpoint, min_salary, max_salary,
			is_benchmark, market_salary, search_keyword,
			created_at, updated_at
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
		SELECT 
			j.id, j.code, j.name, j.description, j.department_id,
			j.e_score, j.c_score, j.r_score, j.we_weight, j.wc_weight, j.wr_weight,
			j.salary_spread, j.job_score, j.midpoint, j.min_salary, j.max_salary,
			j.is_benchmark, j.market_salary, j.search_keyword,
			j.created_at, j.updated_at
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
	query := `
		SELECT 
			id, code, name, description, department_id,
			e_score, c_score, r_score, we_weight, wc_weight, wr_weight,
			salary_spread, job_score, midpoint, min_salary, max_salary,
			is_benchmark, market_salary, search_keyword,
			created_at, updated_at 
		FROM job_descriptions 
		WHERE id = $1
	`
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
		SET 
			name = :name, 
			description = :description, 
			department_id = :department_id,
			e_score = :e_score,
			c_score = :c_score,
			r_score = :r_score,
			we_weight = :we_weight,
			wc_weight = :wc_weight,
			wr_weight = :wr_weight,
			salary_spread = :salary_spread,
			job_score = :job_score,
			midpoint = :midpoint,
			min_salary = :min_salary,
			max_salary = :max_salary,
			is_benchmark = :is_benchmark,
			market_salary = :market_salary,
			search_keyword = :search_keyword,
			updated_at = NOW()
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

func (r *PayrollRepository) GetBenchmarkJobPositions(ctx context.Context) ([]entity.JobPosition, error) {
	query := `
		SELECT 
			id, code, name, description, department_id,
			e_score, c_score, r_score, we_weight, wc_weight, wr_weight,
			salary_spread, job_score, midpoint, min_salary, max_salary,
			is_benchmark, market_salary, search_keyword,
			created_at, updated_at 
		FROM job_descriptions 
		WHERE is_benchmark = TRUE AND job_score > 0 AND market_salary > 0
		ORDER BY name ASC
	`
	var list []entity.JobPosition
	err := r.DB.SelectContext(ctx, &list, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.JobPosition{}, nil
		}
		logger.Error("PayrollRepository:GetBenchmarkJobPositions:Error %v", err)
		return nil, err
	}
	return list, nil
}

func (r *PayrollRepository) GetAllJobPositions(ctx context.Context) ([]entity.JobPosition, error) {
	query := `
		SELECT 
			id, code, name, description, department_id,
			e_score, c_score, r_score, we_weight, wc_weight, wr_weight,
			salary_spread, job_score, midpoint, min_salary, max_salary,
			is_benchmark, market_salary, search_keyword,
			created_at, updated_at 
		FROM job_descriptions
		ORDER BY name ASC
	`
	var list []entity.JobPosition
	err := r.DB.SelectContext(ctx, &list, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.JobPosition{}, nil
		}
		logger.Error("PayrollRepository:GetAllJobPositions:Error %v", err)
		return nil, err
	}
	return list, nil
}
