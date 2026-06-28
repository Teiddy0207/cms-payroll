package mapper

import (
	"cal-salary/modules/payroll/dto"
	"cal-salary/modules/payroll/entity"

	"github.com/google/uuid"
)

func ToEvaluationEntity(req *dto.BatchEvaluationRequest, input dto.CompetencyScoreInput, evaluatorID uuid.UUID) *entity.CompetencyEvaluation {
	return &entity.CompetencyEvaluation{
		EmployeeID:       req.EmployeeID,
		EvaluatorID:      evaluatorID,
		EvaluationPeriod: req.EvaluationPeriod,
		CompetencyID:     input.CompetencyID,
		Score:            input.Score,
		Weight:           input.Weight,
	}
}

func ToEvaluationDTO(ce *entity.CompetencyEvaluation) *dto.EvaluationResponse {
	return &dto.EvaluationResponse{
		ID:               ce.ID,
		EmployeeID:       ce.EmployeeID,
		EvaluatorID:      ce.EvaluatorID,
		EvaluationPeriod: ce.EvaluationPeriod,
		CompetencyID:     ce.CompetencyID,
		Score:            ce.Score,
		Weight:           ce.Weight,
		CreatedAt:        ce.CreatedAt,
		UpdatedAt:        ce.UpdatedAt,
	}
}

func ToPaginatedEvaluationDTO(items []entity.CompetencyEvaluation, totalItems, pageNumber, pageSize int) *dto.PaginatedEvaluationDTO {
	responses := make([]dto.EvaluationResponse, len(items))
	for i, ce := range items {
		responses[i] = *ToEvaluationDTO(&ce)
	}

	totalPages := 0
	if pageSize > 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
	}

	return &dto.PaginatedEvaluationDTO{
		Items:      responses,
		TotalItems: totalItems,
		TotalPages: totalPages,
		PageNumber: pageNumber,
		PageSize:   pageSize,
	}
}
