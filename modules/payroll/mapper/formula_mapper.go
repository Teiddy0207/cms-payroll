package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToFormulaEntity(req *dto.CreateFormulaRequest) *entity.PayrollFormula {
	return &entity.PayrollFormula{
		VariableName: req.VariableName,
		Expression:   req.Expression,
		StartDate:    req.StartDate,
		EndDate:      req.EndDate,
		Description:  req.Description,
	}
}

func ToFormulaDTO(formula *entity.PayrollFormula) *dto.FormulaResponse {
	return &dto.FormulaResponse{
		ID:           formula.ID,
		VariableName: formula.VariableName,
		Expression:   formula.Expression,
		StartDate:    formula.StartDate,
		EndDate:      formula.EndDate,
		Description:  formula.Description,
		CreatedAt:    formula.CreatedAt,
		UpdatedAt:    formula.UpdatedAt,
	}
}

func ToFormulaDTOList(formulas []entity.PayrollFormula) []dto.FormulaResponse {
	list := make([]dto.FormulaResponse, len(formulas))
	for i, f := range formulas {
		list[i] = *ToFormulaDTO(&f)
	}
	return list
}
