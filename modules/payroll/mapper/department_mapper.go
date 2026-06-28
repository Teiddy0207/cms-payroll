package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToDepartmentEntity(req *dto.CreateDepartmentRequest) *entity.Department {
	return &entity.Department{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		ManagerID:   req.ManagerID,
	}
}

func ToDepartmentDTO(dept *entity.Department) *dto.DepartmentResponse {
	return &dto.DepartmentResponse{
		ID:          dept.ID,
		Code:        dept.Code,
		Name:        dept.Name,
		Description: dept.Description,
		ParentID:    dept.ParentID,
		ManagerID:   dept.ManagerID,
		CreatedAt:   dept.CreatedAt,
		UpdatedAt:   dept.UpdatedAt,
	}
}

func ToPaginatedDepartmentDTO(items []entity.Department, totalItems, pageNumber, pageSize int) *dto.PaginatedDepartmentDTO {
	responses := make([]dto.DepartmentResponse, len(items))
	for i, d := range items {
		responses[i] = *ToDepartmentDTO(&d)
	}

	return &dto.PaginatedDepartmentDTO{
		Items:      responses,
		TotalItems: totalItems,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}
