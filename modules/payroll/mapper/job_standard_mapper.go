package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToJobStandardEntity(req *dto.CreateJobStandardRequest) *entity.JobStandard {
	return &entity.JobStandard{
		StandardCode:   req.StandardCode,
		Name:           req.Name,
		AllowanceValue: req.AllowanceValue,
		Description:    req.Description,
	}
}

func ToJobStandardDTO(js *entity.JobStandard) *dto.JobStandardResponse {
	return &dto.JobStandardResponse{
		ID:             js.ID,
		StandardCode:   js.StandardCode,
		Name:           js.Name,
		AllowanceValue: js.AllowanceValue,
		Description:    js.Description,
		CreatedAt:      js.CreatedAt,
		UpdatedAt:      js.UpdatedAt,
	}
}

func ToPaginatedJobStandardDTO(items []entity.JobStandard, totalItems, pageNumber, pageSize int) *dto.PaginatedJobStandardDTO {
	responses := make([]dto.JobStandardResponse, len(items))
	for i, js := range items {
		responses[i] = *ToJobStandardDTO(&js)
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}

	return &dto.PaginatedJobStandardDTO{
		Items:      responses,
		TotalItems: totalItems,
		TotalPages: totalPages,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}
