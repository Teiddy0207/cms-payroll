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
	EScore       float64    `json:"e_score"`
	CScore       float64    `json:"c_score"`
	RScore       float64    `json:"r_score"`
	WEWeight     float64    `json:"we_weight"`
	WCWeight     float64    `json:"wc_weight"`
	WRWeight     float64    `json:"wr_weight"`
	SalarySpread float64    `json:"salary_spread"`
}

type UpdateJobPositionRequest struct {
	Name         string     `json:"name" validate:"required,min=2,max=255"`
	Description  *string    `json:"description,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	EScore       float64    `json:"e_score"`
	CScore       float64    `json:"c_score"`
	RScore       float64    `json:"r_score"`
	WEWeight     float64    `json:"we_weight"`
	WCWeight     float64    `json:"wc_weight"`
	WRWeight     float64    `json:"wr_weight"`
	SalarySpread float64    `json:"salary_spread"`
}

type JobPositionResponse struct {
	ID           uuid.UUID  `json:"id"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	Description  *string    `json:"description,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	EScore       float64    `json:"e_score"`
	CScore       float64    `json:"c_score"`
	RScore       float64    `json:"r_score"`
	WEWeight     float64    `json:"we_weight"`
	WCWeight     float64    `json:"wc_weight"`
	WRWeight     float64    `json:"wr_weight"`
	SalarySpread float64    `json:"salary_spread"`
	JobScore     float64    `json:"job_score"`
	Midpoint     float64    `json:"midpoint"`
	MinSalary    float64    `json:"min_salary"`
	MaxSalary    float64    `json:"max_salary"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}


type PaginatedJobPositionDTO = dto.Pagination[JobPositionResponse]
