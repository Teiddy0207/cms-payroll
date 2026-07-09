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

func (s *PayrollService) CreatePayrollFormula(ctx context.Context, req *dto.CreateFormulaRequest) (*dto.FormulaResponse, *errors.AppError) {
	if req == nil {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "request is required", nil)
	}

	if err := s.validateFormulaExpression(ctx, req.Expression); err != nil {
		return nil, errors.NewAppError(errors.ErrInvalidInput, "Cú pháp công thức không hợp lệ: "+err.Error(), err)
	}

	formula := mapper.ToFormulaEntity(req)
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
		VariableName: req.VariableName,
		Expression:   req.Expression,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Description:  req.Description,
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
