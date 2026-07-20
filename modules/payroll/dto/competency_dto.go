package dto

import (
	coredto "cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

type CreateCompetencyRequest struct {
	Code        string     `json:"code" validate:"required,min=2,max=100"`
	Name        string     `json:"name" validate:"required,min=2,max=255"`
	PointValue  int        `json:"point_value" validate:"required,min=0"`
	Description *string    `json:"description,omitempty"`
	TypeID      *uuid.UUID `json:"type_id,omitempty"`
}

type UpdateCompetencyRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=255"`
	PointValue  int        `json:"point_value" validate:"required,min=0"`
	Description *string    `json:"description,omitempty"`
	TypeID      *uuid.UUID `json:"type_id,omitempty"`
}

type CompetencyResponse struct {
	ID          uuid.UUID  `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	PointValue  int        `json:"point_value"`
	Description *string    `json:"description,omitempty"`
	TypeID      *uuid.UUID `json:"type_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PaginatedCompetencyDTO = coredto.Pagination[CompetencyResponse]
