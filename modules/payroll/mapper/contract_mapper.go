package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
	"time"
)

func ToContractEntity(req *dto.CreateContractRequest) *entity.Contract {
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	var endDate *time.Time
	if req.EndDate != nil && *req.EndDate != "" {
		parsed, _ := time.Parse("2006-01-02", *req.EndDate)
		endDate = &parsed
	}

	return &entity.Contract{
		EmployeeID:       req.EmployeeID,
		ContractCode:     req.ContractCode,
		PositionBaseRate: req.PositionBaseRate,
		StartDate:        startDate,
		EndDate:          endDate,
		Status:           req.Status,
	}
}

func ToContractDTO(c *entity.Contract) *dto.ContractResponse {
	return &dto.ContractResponse{
		ID:               c.ID,
		EmployeeID:       c.EmployeeID,
		ContractCode:     c.ContractCode,
		PositionBaseRate: c.PositionBaseRate,
		StartDate:        c.StartDate,
		EndDate:          c.EndDate,
		Status:           c.Status,
		CreatedAt:        c.CreatedAt,
		UpdatedAt:        c.UpdatedAt,
	}
}

func ToPaginatedContractDTO(items []entity.Contract, totalItems, pageNumber, pageSize int) *dto.PaginatedContractDTO {
	responses := make([]dto.ContractResponse, len(items))
	for i, c := range items {
		responses[i] = *ToContractDTO(&c)
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}

	return &dto.PaginatedContractDTO{
		Items:      responses,
		TotalItems: totalItems,
		TotalPages: totalPages,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}
