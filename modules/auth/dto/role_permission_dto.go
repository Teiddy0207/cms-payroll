package dto

import (
	"cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

type RolePermissionRequest struct {
	dto.BaseRequest
	RoleID       uuid.UUID   `json:"role_id"`
	PermissionID []uuid.UUID `json:"permission_id"`
	GrantedBy    uuid.UUID   `json:"-"`
}

type RolePermissionResponse struct {
	ID           uuid.UUID   `json:"id"`
	RoleID       uuid.UUID   `json:"role_id"`
	PermissionID []uuid.UUID `json:"permission_id"`
	GrantedBy    *uuid.UUID  `json:"granted_by"`
	GrantedAt    time.Time   `json:"granted_at"`
}

type RolePermissionsDTO struct {
	Permissions []PermissionResponse `json:"permissions"`
}

type PaginatedRolePermissionDTO = dto.Pagination[RolePermissionResponse]
