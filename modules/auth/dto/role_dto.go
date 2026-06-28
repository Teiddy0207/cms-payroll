package dto

import (
	"cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

type RoleRequest struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	IsSystem    bool    `json:"is_system"`
	IsActive    bool    `json:"is_active"`
}

type RoleResponse struct {
	ID          uuid.UUID            `json:"id"`
	Name        string               `json:"name"`
	Slug        string               `json:"slug"`
	Description *string              `json:"description"`
	IsSystem    bool                 `json:"is_system"`
	IsActive    bool                 `json:"is_active"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}

type UserRoleRequest struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	UserID      uuid.UUID `json:"user_id"`
	Description *string   `json:"description,omitempty"`
}

type PaginatedRoleDTO = dto.Pagination[RoleResponse]
