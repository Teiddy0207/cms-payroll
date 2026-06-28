package dto

import (
	"cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

type CreateContractRequest struct {
	EmployeeID       uuid.UUID `json:"employee_id" validate:"required"`
	ContractCode     string    `json:"contract_code" validate:"required,min=2,max=100"`
	PositionBaseRate float64   `json:"position_base_rate" validate:"required,min=0"`
	StartDate        string    `json:"start_date" validate:"required"` // format: YYYY-MM-DD
	EndDate          *string   `json:"end_date,omitempty"`             // format: YYYY-MM-DD
	Status           string    `json:"status" validate:"required,oneof=ACTIVE EXPIRED TERMINATED"`
}

type UpdateContractRequest struct {
	PositionBaseRate float64   `json:"position_base_rate" validate:"required,min=0"`
	StartDate        string    `json:"start_date" validate:"required"`
	EndDate          *string   `json:"end_date,omitempty"`
	Status           string    `json:"status" validate:"required,oneof=ACTIVE EXPIRED TERMINATED"`
}

type ContractResponse struct {
	ID               uuid.UUID  `json:"id"`
	EmployeeID       uuid.UUID  `json:"employee_id"`
	ContractCode     string     `json:"contract_code"`
	PositionBaseRate float64    `json:"position_base_rate"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type PaginatedContractDTO = dto.Pagination[ContractResponse]
