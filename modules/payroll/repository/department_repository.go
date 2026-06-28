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
// Departments Repository Implementation
// ==========================================

func (r *PayrollRepository) CreateDepartment(ctx context.Context, dept *entity.Department) (*entity.Department, error) {
	if dept.ID == uuid.Nil {
		dept.ID = uuid.New()
	}
	query := `
		INSERT INTO departments (id, code, name, description, parent_id, manager_id, created_at, updated_at)
		VALUES (:id, :code, :name, :description, :parent_id, :manager_id, NOW(), NOW())
		RETURNING id, code, name, description, parent_id, manager_id, created_at, updated_at
	`
	rows, err := r.DB.NamedQueryContext(ctx, query, dept)
	if err != nil {
		logger.Error("PayrollRepository:CreateDepartment:Error %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var createdDept entity.Department
		if err := rows.StructScan(&createdDept); err != nil {
			logger.Error("PayrollRepository:CreateDepartment:Scan:Error %v", err)
			return nil, err
		}
		return &createdDept, nil
	}
	return nil, errors.New("failed to create department")
}

func (r *PayrollRepository) GetDepartments(ctx context.Context, qp params.QueryParams) ([]entity.Department, int, error) {
	offset := (qp.PageNumber - 1) * qp.PageSize
	baseQuery := `FROM departments d`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if qp.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(d.name ILIKE $%d OR d.code ILIKE $%d)", argIndex, argIndex))
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
		logger.Error("PayrollRepository:GetDepartments - Count", err)
		return nil, 0, err
	}

	dataQuery := `
		SELECT d.id, d.code, d.name, d.description, d.parent_id, d.manager_id, d.created_at, d.updated_at
	` + baseQuery + whereClause + ` ORDER BY d.name ASC, d.created_at DESC`

	// Add limit and offset
	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, qp.PageSize, offset)

	var depts []entity.Department
	err = r.DB.SelectContext(ctx, &depts, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.Department{}, 0, nil
		}
		logger.Error("PayrollRepository:GetDepartments - Select", err)
		return nil, 0, err
	}

	return depts, totalItems, nil
}

func (r *PayrollRepository) GetDepartmentByID(ctx context.Context, id uuid.UUID) (*entity.Department, error) {
	var dept entity.Department
	query := `SELECT * FROM departments WHERE id = $1`
	err := r.DB.GetContext(ctx, &dept, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("PayrollRepository:GetDepartmentByID - Select", err)
		return nil, err
	}
	return &dept, nil
}

func (r *PayrollRepository) UpdateDepartment(ctx context.Context, id uuid.UUID, dept *entity.Department) error {
	query := `
		UPDATE departments
		SET name = :name, description = :description, parent_id = :parent_id, manager_id = :manager_id, updated_at = NOW()
		WHERE id = :id
	`
	dept.ID = id
	result, err := r.DB.NamedExecContext(ctx, query, dept)
	if err != nil {
		logger.Error("PayrollRepository:UpdateDepartment:Error %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("department not found")
	}
	return nil
}

func (r *PayrollRepository) DeleteDepartment(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM departments WHERE id = $1`
	_, err := r.DB.SQLx().ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("PayrollRepository:DeleteDepartment:Error %v", err)
		return err
	}
	return nil
}
