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
// Contracts Repository Implementation
// ==========================================

func (r *PayrollRepository) CreateContract(ctx context.Context, contract *entity.Contract) (*entity.Contract, error) {
	if contract.ID == uuid.Nil {
		contract.ID = uuid.New()
	}
	query := `
		INSERT INTO contracts (id, employee_id, contract_code, position_base_rate, start_date, end_date, status, created_at, updated_at)
		VALUES (:id, :employee_id, :contract_code, :position_base_rate, :start_date, :end_date, :status, NOW(), NOW())
		RETURNING id, employee_id, contract_code, position_base_rate, start_date, end_date, status, created_at, updated_at
	`
	rows, err := r.DB.NamedQueryContext(ctx, query, contract)
	if err != nil {
		logger.Error("PayrollRepository:CreateContract:Error %v", err)
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		var created entity.Contract
		if err := rows.StructScan(&created); err != nil {
			logger.Error("PayrollRepository:CreateContract:Scan:Error %v", err)
			return nil, err
		}
		return &created, nil
	}
	return nil, errors.New("failed to create contract")
}

func (r *PayrollRepository) GetContracts(ctx context.Context, qp params.QueryParams) ([]entity.Contract, int, error) {
	offset := (qp.PageNumber - 1) * qp.PageSize
	baseQuery := `FROM contracts c`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if qp.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(c.contract_code ILIKE $%d)", argIndex))
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
		logger.Error("PayrollRepository:GetContracts - Count", err)
		return nil, 0, err
	}

	dataQuery := `
		SELECT c.id, c.employee_id, c.contract_code, c.position_base_rate, c.start_date, c.end_date, c.status, c.created_at, c.updated_at
	` + baseQuery + whereClause + ` ORDER BY c.created_at DESC`

	dataQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, qp.PageSize, offset)

	var contracts []entity.Contract
	err = r.DB.SelectContext(ctx, &contracts, dataQuery, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.Contract{}, 0, nil
		}
		logger.Error("PayrollRepository:GetContracts - Select", err)
		return nil, 0, err
	}

	return contracts, totalItems, nil
}

func (r *PayrollRepository) GetContractByID(ctx context.Context, id uuid.UUID) (*entity.Contract, error) {
	var contract entity.Contract
	query := `SELECT * FROM contracts WHERE id = $1`
	err := r.DB.GetContext(ctx, &contract, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logger.Error("PayrollRepository:GetContractByID - Select", err)
		return nil, err
	}
	return &contract, nil
}

func (r *PayrollRepository) UpdateContract(ctx context.Context, id uuid.UUID, contract *entity.Contract) error {
	query := `
		UPDATE contracts
		SET position_base_rate = :position_base_rate, start_date = :start_date, end_date = :end_date, status = :status, updated_at = NOW()
		WHERE id = :id
	`
	contract.ID = id
	result, err := r.DB.NamedExecContext(ctx, query, contract)
	if err != nil {
		logger.Error("PayrollRepository:UpdateContract:Error %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("contract not found")
	}
	return nil
}

func (r *PayrollRepository) DeleteContract(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM contracts WHERE id = $1`
	_, err := r.DB.SQLx().ExecContext(ctx, query, id)
	if err != nil {
		logger.Error("PayrollRepository:DeleteContract:Error %v", err)
		return err
	}
	return nil
}

func (r *PayrollRepository) HasActiveContract(ctx context.Context, employeeID uuid.UUID, excludeContractID *uuid.UUID) (bool, error) {
	var count int
	var err error
	if excludeContractID != nil {
		query := `SELECT COUNT(*) FROM contracts WHERE employee_id = $1 AND status = 'ACTIVE' AND id <> $2`
		err = r.DB.GetContext(ctx, &count, query, employeeID, *excludeContractID)
	} else {
		query := `SELECT COUNT(*) FROM contracts WHERE employee_id = $1 AND status = 'ACTIVE'`
		err = r.DB.GetContext(ctx, &count, query, employeeID)
	}
	if err != nil {
		logger.Error("PayrollRepository:HasActiveContract:Error %v", err)
		return false, err
	}
	return count > 0, nil
}
