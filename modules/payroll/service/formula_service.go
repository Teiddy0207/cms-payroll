package service

import (
	"cal-salary/core/errors"
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
	"cal-salary/modules/payroll/mapper"
	"context"

	"github.com/expr-lang/expr"
	"github.com/google/uuid"
)

func (s *PayrollService) validateFormulaExpression(ctx context.Context, expression string) error {
	env := map[string]interface{}{
		"P1":           0.0,
		"p1":           0.0,
		"P2":           0.0,
		"p2":           0.0,
		"P3":           0.0,
		"p3":           0.0,
		"GROSS_SALARY": 0.0,
		"gross_salary": 0.0,
		"TAX":          0.0,
		"tax":          0.0,
		"NET_SALARY":   0.0,
		"net_salary":   0.0,
	}

	formulas, err := s.repo.GetPayrollFormulas(ctx)
	if err == nil {
		for _, f := range formulas {
			env[f.VariableName] = 0.0
		}
	}

	_, err = expr.Compile(expression, expr.Env(env))
	return err
}

// checkNoCycle re-runs the topological sort across all existing formulas plus the
// candidate being created/updated, rejecting the save if it introduces a circular
// dependency. excludeID identifies the formula being edited (uuid.Nil for create).
func (s *PayrollService) checkNoCycle(ctx context.Context, candidate entity.PayrollFormula, excludeID uuid.UUID) error {
	formulas, err := s.repo.GetPayrollFormulas(ctx)
	if err != nil {
		return nil
	}

	set := make([]entity.PayrollFormula, 0, len(formulas)+1)
	for _, f := range formulas {
		if f.ID != excludeID {
			set = append(set, f)
		}
	}
	set = append(set, candidate)

	if _, err := sortFormulas(set); err != nil {
		return err
	}
	return nil
}

func (s *PayrollService) CreatePayrollFormula(ctx context.Context, req *dto.CreateFormulaRequest) (*dto.FormulaResponse, *errors.AppError) {
	if req == nil {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "request is required", nil)
	}

	if err := s.validateFormulaExpression(ctx, req.Expression); err != nil {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "Cú pháp công thức không hợp lệ: "+err.Error(), err)
	}

	formula := mapper.ToFormulaEntity(req)

	if err := s.checkNoCycle(ctx, *formula, uuid.Nil); err != nil {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "Công thức tạo ra vòng lặp phụ thuộc: "+err.Error(), err)
	}

	created, err := s.repo.CreatePayrollFormula(ctx, formula)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to create payroll formula", err)
	}

	return mapper.ToFormulaDTO(created), nil
}

func (s *PayrollService) GetPayrollFormulas(ctx context.Context) ([]dto.FormulaResponse, *errors.AppError) {
	formulas, err := s.repo.GetPayrollFormulas(ctx)
	if err != nil {
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to fetch payroll formulas", err)
	}
	return mapper.ToFormulaDTOList(formulas), nil
}

func (s *PayrollService) UpdatePayrollFormula(ctx context.Context, id uuid.UUID, req *dto.UpdateFormulaRequest) *errors.AppError {
	if req == nil {
		return errors.NewAppError(errors.ErrInvalidInput, "request is required", nil)
	}

	if err := s.validateFormulaExpression(ctx, req.Expression); err != nil {
		return errors.NewAppError(errors.ErrInvalidInput, "Cú pháp công thức không hợp lệ: "+err.Error(), err)
	}

	formula := &entity.PayrollFormula{
		ID:           id,
		VariableName: req.VariableName,
		Expression:   req.Expression,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Description:  req.Description,
	}

	if err := s.checkNoCycle(ctx, *formula, id); err != nil {
		return errors.NewAppError(errors.ErrInvalidInput, "Công thức tạo ra vòng lặp phụ thuộc: "+err.Error(), err)
	}

	err := s.repo.UpdatePayrollFormula(ctx, id, formula)
	if err != nil {
		return errors.NewAppError(errors.ErrInternalServer, "failed to update payroll formula", err)
	}
	return nil
}

func (s *PayrollService) DeletePayrollFormula(ctx context.Context, id uuid.UUID) *errors.AppError {
	err := s.repo.DeletePayrollFormula(ctx, id)
	if err != nil {
		return errors.NewAppError(errors.ErrInternalServer, "failed to delete payroll formula", err)
	}
	return nil
}
