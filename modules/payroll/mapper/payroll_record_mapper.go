package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToSalaryCalculationItem(record *entity.PayrollRecord, profile *entity.UserProfile, details []entity.PayrollRecordDetail, period string) *dto.SalaryCalculationItem {
	var breakdown map[string]any
	if len(details) > 0 {
		breakdown = make(map[string]any)
		for _, d := range details {
			breakdown[d.Description] = d.Amount
		}
	}

	preview := &dto.SalaryPreviewResponse{
		EmployeeID:   record.EmployeeID,
		FullName:     profile.FullName,
		Period:       period,
		P1Total:      record.P1Value,
		P2Total:      record.P2Value,
		SubtotalP1P2: record.GrossSalary,
		P1:           record.P1Value,
		P2:           record.P2Value,
		Total:        record.GrossSalary,
		SalaryP1:     record.P1Value,
		SalaryP2:     record.P2Value,
		TotalSalary:  record.GrossSalary,
		Breakdown:    breakdown,
	}

	return &dto.SalaryCalculationItem{
		EmployeeID:   record.EmployeeID,
		FullName:     profile.FullName,
		DepartmentID: profile.DepartmentID,
		PositionID:   profile.PositionID,
		Status:       record.Status,
		Preview:      preview,
	}
}
