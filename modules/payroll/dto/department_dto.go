package dto

import (
	"cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

type CreateDepartmentRequest struct {
	Code        string     `json:"code" validate:"required,min=2,max=50"`
	Name        string     `json:"name" validate:"required,min=2,max=255"`
	Description *string    `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	ManagerID   *uuid.UUID `json:"manager_id,omitempty"`
}

type UpdateDepartmentRequest struct {
	Name        string     `json:"name" validate:"required,min=2,max=255"`
	Description *string    `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	ManagerID   *uuid.UUID `json:"manager_id,omitempty"`
}

type DepartmentResponse struct {
	ID          uuid.UUID  `json:"id"`
	Code        string     `json:"code"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	ManagerID   *uuid.UUID `json:"manager_id,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PaginatedDepartmentDTO = dto.Pagination[DepartmentResponse]
