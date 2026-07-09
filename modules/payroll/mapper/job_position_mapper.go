package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToJobPositionEntity(req *dto.CreateJobPositionRequest) *entity.JobPosition {
	return &entity.JobPosition{
		Code:         req.Code,
		Name:         req.Name,
		Description:  req.Description,
		DepartmentID: req.DepartmentID,
		EScore:       req.EScore,
		CScore:       req.CScore,
		RScore:       req.RScore,
		WEWeight:     req.WEWeight,
		WCWeight:     req.WCWeight,
		WRWeight:     req.WRWeight,
		SalarySpread: req.SalarySpread,
	}
}

func ToJobPositionDTO(pos *entity.JobPosition) *dto.JobPositionResponse {
	return &dto.JobPositionResponse{
		ID:           pos.ID,
		Code:         pos.Code,
		Name:         pos.Name,
		Description:  pos.Description,
		DepartmentID: pos.DepartmentID,
		EScore:       pos.EScore,
		CScore:       pos.CScore,
		RScore:       pos.RScore,
		WEWeight:     pos.WEWeight,
		WCWeight:     pos.WCWeight,
		WRWeight:     pos.WRWeight,
		SalarySpread: pos.SalarySpread,
		JobScore:     pos.JobScore,
		Midpoint:     pos.Midpoint,
		MinSalary:    pos.MinSalary,
		MaxSalary:    pos.MaxSalary,
		CreatedAt:    pos.CreatedAt,
		UpdatedAt:    pos.UpdatedAt,
	}
}

func ToPaginatedJobPositionDTO(items []entity.JobPosition, totalItems, pageNumber, pageSize int) *dto.PaginatedJobPositionDTO {
	responses := make([]dto.JobPositionResponse, len(items))
	for i, p := range items {
		responses[i] = *ToJobPositionDTO(&p)
	}

	return &dto.PaginatedJobPositionDTO{
		Items:      responses,
		TotalItems: totalItems,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}
