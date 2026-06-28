package dto

import (
	"cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

type CreateJobStandardRequest struct {
	StandardCode   string  `json:"standard_code" validate:"required,min=2,max=100"`
	Name           string  `json:"name" validate:"required,min=2,max=255"`
	AllowanceValue float64 `json:"allowance_value" validate:"required,min=0"`
	Description    *string `json:"description,omitempty"`
}

type UpdateJobStandardRequest struct {
	Name           string  `json:"name" validate:"required,min=2,max=255"`
	AllowanceValue float64 `json:"allowance_value" validate:"required,min=0"`
	Description    *string `json:"description,omitempty"`
}

type JobStandardResponse struct {
	ID             uuid.UUID `json:"id"`
	StandardCode   string    `json:"standard_code"`
	Name           string    `json:"name"`
	AllowanceValue float64   `json:"allowance_value"`
	Description    *string   `json:"description,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type PaginatedJobStandardDTO = dto.Pagination[JobStandardResponse]
