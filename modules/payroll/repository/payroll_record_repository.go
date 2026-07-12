package repository

import (
	"cal-salary/modules/payroll/entity"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

func (r *PayrollRepository) GetOrCreatePeriod(ctx context.Context, month, year int) (*entity.PayrollPeriod, error) {
	var period entity.PayrollPeriod
	query := `SELECT id, month, year, status, locked_at, disbursed_at, created_at, updated_at FROM payroll_periods WHERE month = $1 AND year = $2`
	err := r.DB.SQLx().GetContext(ctx, &period, query, month, year)
	if err == nil {
		return &period, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}

	newID := uuid.New()
	insertQuery := `
		INSERT INTO payroll_periods (id, month, year, status, created_at, updated_at)
		VALUES ($1, $2, $3, 'DRAFT', NOW(), NOW())
		RETURNING id, month, year, status, locked_at, disbursed_at, created_at, updated_at
	`
	err = r.DB.SQLx().GetContext(ctx, &period, insertQuery, newID, month, year)
	if err != nil {
		return nil, err
	}
	return &period, nil
}

func (r *PayrollRepository) GetPayrollPeriodByMonthYear(ctx context.Context, month, year int) (*entity.PayrollPeriod, error) {
	var period entity.PayrollPeriod
	query := `SELECT id, month, year, status, locked_at, disbursed_at, created_at, updated_at FROM payroll_periods WHERE month = $1 AND year = $2`
	err := r.DB.SQLx().GetContext(ctx, &period, query, month, year)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &period, nil
}

func (r *PayrollRepository) UpsertPayrollRecord(ctx context.Context, record *entity.PayrollRecord, details []entity.PayrollRecordDetail) error {
	tx, err := r.DB.SQLx().BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if record.ID == uuid.Nil {
		record.ID = uuid.New()
	}

	recordQuery := `
		INSERT INTO payroll_records (id, period_id, employee_id, p1_value, p2_value, p3_value, gross_salary, tax, net_salary, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		ON CONFLICT (period_id, employee_id)
		DO UPDATE SET p1_value = EXCLUDED.p1_value, p2_value = EXCLUDED.p2_value, p3_value = EXCLUDED.p3_value, gross_salary = EXCLUDED.gross_salary, tax = EXCLUDED.tax, net_salary = EXCLUDED.net_salary, status = EXCLUDED.status, updated_at = NOW()
		RETURNING id
	`
	var actualID uuid.UUID
	err = tx.GetContext(ctx, &actualID, recordQuery,
		record.ID, record.PeriodID, record.EmployeeID,
		record.P1Value, record.P2Value, record.P3Value,
		record.GrossSalary, record.Tax, record.NetSalary,
		record.Status,
	)
	if err != nil {
		return err
	}
	record.ID = actualID

	_, err = tx.ExecContext(ctx, `DELETE FROM payroll_record_details WHERE record_id = $1`, record.ID)
	if err != nil {
		return err
	}

	detailQuery := `
		INSERT INTO payroll_record_details (id, record_id, component, description, source, amount, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	for _, d := range details {
		detailID := uuid.New()
		_, err = tx.ExecContext(ctx, detailQuery, detailID, record.ID, d.Component, d.Description, d.Source, d.Amount)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PayrollRepository) GetPayrollRecords(ctx context.Context, periodID uuid.UUID) ([]entity.PayrollRecord, error) {
	var list []entity.PayrollRecord
	query := `SELECT id, period_id, employee_id, p1_value, p2_value, p3_value, gross_salary, tax, net_salary, status, created_at, updated_at FROM payroll_records WHERE period_id = $1`
	err := r.DB.SQLx().SelectContext(ctx, &list, query, periodID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.PayrollRecord{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (r *PayrollRepository) GetPayrollRecordDetails(ctx context.Context, recordID uuid.UUID, departmentID *uuid.UUID) ([]entity.PayrollRecordDetail, error) {
	var list []entity.PayrollRecordDetail
	var err error
	if departmentID != nil {
		query := `
			SELECT d.id, d.record_id, d.component, d.description, d.source, d.amount, d.created_at 
			FROM payroll_record_details d
			JOIN payroll_records pr ON d.record_id = pr.id
			JOIN user_profiles u ON pr.employee_id = u.id
			WHERE d.record_id = $1 AND u.department_id = $2
		`
		err = r.DB.SQLx().SelectContext(ctx, &list, query, recordID, *departmentID)
	} else {
		query := `SELECT id, record_id, component, description, source, amount, created_at FROM payroll_record_details WHERE record_id = $1`
		err = r.DB.SQLx().SelectContext(ctx, &list, query, recordID)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.PayrollRecordDetail{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (r *PayrollRepository) GetPayrollFormulasByPeriod(ctx context.Context, start, end time.Time) ([]entity.PayrollFormula, error) {
	var list []entity.PayrollFormula
	query := `
		SELECT id, variable_name, expression, start_date, end_date, description, created_at, updated_at
		FROM payroll_formulas
		WHERE start_date <= $1 AND (end_date IS NULL OR end_date >= $2)
	`
	err := r.DB.SQLx().SelectContext(ctx, &list, query, end, start)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.PayrollFormula{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (r *PayrollRepository) CreatePayrollFormula(ctx context.Context, formula *entity.PayrollFormula) (*entity.PayrollFormula, error) {
	if formula.ID == uuid.Nil {
		formula.ID = uuid.New()
	}
	query := `
		INSERT INTO payroll_formulas (id, variable_name, expression, start_date, end_date, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, variable_name, expression, start_date, end_date, description, created_at, updated_at
	`
	var created entity.PayrollFormula
	err := r.DB.SQLx().GetContext(ctx, &created, query,
		formula.ID, formula.VariableName, formula.Expression,
		formula.StartDate, formula.EndDate, formula.Description,
	)
	if err != nil {
		return nil, err
	}
	return &created, nil
}

func (r *PayrollRepository) GetPayrollFormulas(ctx context.Context) ([]entity.PayrollFormula, error) {
	var list []entity.PayrollFormula
	query := `SELECT id, variable_name, expression, start_date, end_date, description, created_at, updated_at FROM payroll_formulas ORDER BY start_date DESC`
	err := r.DB.SQLx().SelectContext(ctx, &list, query)
	if err != nil {
		if err == sql.ErrNoRows {
			return []entity.PayrollFormula{}, nil
		}
		return nil, err
	}
	return list, nil
}

func (r *PayrollRepository) UpdatePayrollFormula(ctx context.Context, id uuid.UUID, formula *entity.PayrollFormula) error {
	query := `
		UPDATE payroll_formulas
		SET variable_name = $1, expression = $2, start_date = $3, end_date = $4, description = $5, updated_at = NOW()
		WHERE id = $6
	`
	_, err := r.DB.SQLx().ExecContext(ctx, query,
		formula.VariableName, formula.Expression,
		formula.StartDate, formula.EndDate, formula.Description,
		id,
	)
	return err
}

func (r *PayrollRepository) DeletePayrollFormula(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM payroll_formulas WHERE id = $1`
	_, err := r.DB.SQLx().ExecContext(ctx, query, id)
	return err
}
