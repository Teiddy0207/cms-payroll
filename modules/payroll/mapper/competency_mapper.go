package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"
)

func ToCompetencyEntity(req *dto.CreateCompetencyRequest) *entity.Competency {
	return &entity.Competency{
		Code:        req.Code,
		Name:        req.Name,
		PointValue:  req.PointValue,
		Description: req.Description,
		TypeID:      req.TypeID,
	}
}

func ToCompetencyDTO(c *entity.Competency) *dto.CompetencyResponse {
	return &dto.CompetencyResponse{
		ID:          c.ID,
		Code:        c.Code,
		Name:        c.Name,
		PointValue:  c.PointValue,
		Description: c.Description,
		TypeID:      c.TypeID,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func ToPaginatedCompetencyDTO(items []entity.Competency, totalItems, pageNumber, pageSize int) *dto.PaginatedCompetencyDTO {
	responses := make([]dto.CompetencyResponse, len(items))
	for i, c := range items {
		responses[i] = *ToCompetencyDTO(&c)
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}

	return &dto.PaginatedCompetencyDTO{
		Items:      responses,
		TotalItems: totalItems,
		TotalPages: totalPages,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}
