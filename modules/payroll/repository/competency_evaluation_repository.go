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
// Competency Evaluations Repository
// ==========================================

func (r *PayrollRepository) CreateEvaluation(ctx context.Context, eval *entity.CompetencyEvaluation) (*entity.CompetencyEvaluation, error) {
	if eval.ID == uuid.Nil {
		eval.ID = uuid.New()
	}
	query := `
		INSERT INTO competency_evaluations (id, employee_id, evaluator_id, evaluation_period, competency_id, score, weight, created_at, updated_at)
		VALUES (:id, :employee_id, :evaluator_id, :evaluation_period, :competency_id, :score, :weight, NOW(), NOW())
		RETURNING id, employee_id, evaluator_id, evaluation_period, competency_id, score, weight, created_at, updated_at
	`
	rows, err := r.DB.NamedQueryContext(ctx, query, eval)
	if err != nil {
		logger.Error("PayrollRepository:CreateEvaluation:Error %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var created entity.CompetencyEvaluation
		if err := rows.StructScan(&created); err != nil {
			logger.Error("PayrollRepository:CreateEvaluation:Scan:Error %v", err)
			return nil, err
		}
		return &created, nil
	}
	return nil, errors.New("failed to create competency evaluation")
}

func (r *PayrollRepository) GetEvaluations(ctx context.Context, qp params.QueryParams) ([]entity.CompetencyEvaluation, int, error) {
	offset := (qp.PageNumber - 1) * qp.PageSize
	baseQuery := `FROM competency_evaluations ce`
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Filter by employee_id if present in filters
	if empIDStr, ok := qp.Filters["employee_id"]; ok && empIDStr != "" {
		empID, err := uuid.Parse(empIDStr)
		if err == nil {
			conditions = append(conditions, fmt.Sprintf("ce.employee_id = $%d", argIndex))
			args = append(args, empID)
			argIndex++
		}
	}

	// Filter by period if present in filters
	if period, ok := qp.Filters["evaluation_period"]; ok && period != "" {
		conditions = append(conditions, fmt.Sprintf("ce.evaluation_period = $%d", argIndex))
		args = append(args, period)
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
		logger.Error("PayrollRepository:GetEvaluations - Count", err)
		return nil, 0, err
	}

	dataQuery := `
		SELECT ce.id, ce.employee_id, ce.evaluator_id, ce.evaluation_period, ce.competency_id, ce.score, ce.weight, ce.created_at, ce.updated_at
	` + baseQuery + whereClause + ` ORDER BY ce.created_at DESC`

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, qp.PageSize, offset)

	var evals []entity.CompetencyEvaluation
	err = r.DB.SelectContext(ctx, &evals, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.CompetencyEvaluation{}, 0, nil
		}
		logger.Error("PayrollRepository:GetEvaluations - Select", err)
		return nil, 0, err
	}

	return evals, totalItems, nil
}

func (r *PayrollRepository) GetEvaluationByID(ctx context.Context, id uuid.UUID) (*entity.CompetencyEvaluation, error) {
	var ce entity.CompetencyEvaluation
	query := `SELECT * FROM competency_evaluations WHERE id = $1`
	err := r.DB.GetContext(ctx, &ce, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("PayrollRepository:GetEvaluationByID - Select", err)
		return nil, err
	}
	return &ce, nil
}

func (r *PayrollRepository) UpdateEvaluation(ctx context.Context, id uuid.UUID, eval *entity.CompetencyEvaluation) error {
	query := `
		UPDATE competency_evaluations
		SET score = :score, weight = :weight, updated_at = NOW()
		WHERE id = :id
	`
	eval.ID = id
	result, err := r.DB.NamedExecContext(ctx, query, eval)
	if err != nil {
		logger.Error("PayrollRepository:UpdateEvaluation:Error %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("evaluation not found")
	}
	return nil
}

func (r *PayrollRepository) DeleteEvaluation(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM competency_evaluations WHERE id = $1`
	_, err := r.DB.SQLx().ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("PayrollRepository:DeleteEvaluation:Error %v", err)
		return err
	}
	return nil
}

func (r *PayrollRepository) DeleteEvaluationsByPeriodAndEmployee(ctx context.Context, employeeID uuid.UUID, period string) error {
	query := `DELETE FROM competency_evaluations WHERE employee_id = $1 AND evaluation_period = $2`
	_, err := r.DB.SQLx().ExecContext(ctx, query, employeeID, period)
	if err != nil {
		logger.Error("PayrollRepository:DeleteEvaluationsByPeriodAndEmployee:Error %v", err)
		return err
	}
	return nil
}

func (r *PayrollRepository) GetManagedDepartmentID(ctx context.Context, managerID uuid.UUID) (uuid.UUID, error) {
	var deptID uuid.UUID
	query := `SELECT id FROM departments WHERE manager_id = $1 LIMIT 1`
	err := r.DB.GetContext(ctx, &deptID, query, managerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, nil
		}
		logger.Error("PayrollRepository:GetManagedDepartmentID:Error %v", err)
		return uuid.Nil, err
	}
	return deptID, nil
}

func (r *PayrollRepository) GetEmployeeDepartmentID(ctx context.Context, employeeID uuid.UUID) (uuid.UUID, error) {
	var deptID uuid.UUID
	// employeeID refers to user_profiles.id
	query := `SELECT department_id FROM user_profiles WHERE id = $1 LIMIT 1`
	err := r.DB.GetContext(ctx, &deptID, query, employeeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return uuid.Nil, nil
		}
		logger.Error("PayrollRepository:GetEmployeeDepartmentID:Error %v", err)
		return uuid.Nil, err
	}
	return deptID, nil
}

func (r *PayrollRepository) IsAdminOrDirector(ctx context.Context, userID uuid.UUID) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) 
		FROM user_roles ur
		JOIN roles r ON ur.role_id = r.id
		WHERE ur.user_id = $1 AND ur.is_active = true AND r.slug IN ('admin', 'director')
	`
	err := r.DB.GetContext(ctx, &count, query, userID)
	if err != nil {
		logger.Error("PayrollRepository:IsAdminOrDirector:Error %v", err)
		return false, err
	}
	return count > 0, nil
}
