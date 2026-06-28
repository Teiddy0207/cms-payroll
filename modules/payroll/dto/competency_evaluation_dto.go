package dto

import (
	"cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

type CompetencyScoreInput struct {
	CompetencyID uuid.UUID `json:"competency_id" validate:"required"`
	Score        int       `json:"score" validate:"required,min=0,max=100"`
	Weight       int       `json:"weight" validate:"required,min=1,max=5"`
}

type BatchEvaluationRequest struct {
	EmployeeID       uuid.UUID              `json:"employee_id" validate:"required"`
	EvaluationPeriod string                 `json:"evaluation_period" validate:"required,min=2,max=50"`
	Scores           []CompetencyScoreInput `json:"scores" validate:"required,dive,required"`
}

type EvaluationResponse struct {
	ID               uuid.UUID `json:"id"`
	EmployeeID       uuid.UUID `json:"employee_id"`
	EvaluatorID      uuid.UUID `json:"evaluator_id"`
	EvaluationPeriod string    `json:"evaluation_period"`
	CompetencyID     uuid.UUID `json:"competency_id"`
	Score            int       `json:"score"`
	Weight           int       `json:"weight"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type PaginatedEvaluationDTO = dto.Pagination[EvaluationResponse]
