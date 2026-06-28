package dto

import (
	"cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

// ==========================================
// JobPosition DTOs
// ==========================================

type CreateJobPositionRequest struct {
	Code         string     `json:"code" validate:"required,min=2,max=50"`
	Name         string     `json:"name" validate:"required,min=2,max=255"`
	Description  *string    `json:"description,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
}

type UpdateJobPositionRequest struct {
	Name         string     `json:"name" validate:"required,min=2,max=255"`
	Description  *string    `json:"description,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
}

type JobPositionResponse struct {
	ID           uuid.UUID  `json:"id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Description  *string    `json:"description,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PaginatedJobPositionDTO = dto.Pagination[JobPositionResponse]
