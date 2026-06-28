package mapper

import (
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/entity"

	"github.com/google/uuid"
)

func ToRoleEntity(role *dto.RoleRequest) *entity.Role {
	return &entity.Role{
		Name:        role.Name,
		Slug:        role.Slug,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		IsActive:    role.IsActive,
	}
}

func ToRoleDTO(role *entity.Role, permissions *[]entity.Permission) *dto.RoleResponse {
	response := &dto.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Slug:        role.Slug,
		Description: role.Description,
		IsSystem:    role.IsSystem,
		IsActive:    role.IsActive,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
	}

	if permissions != nil && len(*permissions) > 0 {
		permissionDTOs := make([]dto.PermissionResponse, len(*permissions))
		for i, permission := range *permissions {
			permissionDTOs[i] = *ToPermissionDTO(&permission)
		}
		response.Permissions = permissionDTOs
	}

	return response
}

func ToPaginatedRoleDTO(entity *entity.PaginatedRoleEntity, rolePermissionsMap map[uuid.UUID]*[]entity.Permission) *dto.PaginatedRoleDTO {

	roleResponses := make([]dto.RoleResponse, len(entity.Items))
	for i, role := range entity.Items {
		permissions := rolePermissionsMap[role.ID]
		roleResponses[i] = *ToRoleDTO(&role, permissions)
	}

	// Tính total pages
	totalPages := 0
	if entity.PageSize > 0 {
		totalPages = (entity.TotalItems + entity.PageSize - 1) / entity.PageSize
	}

	return &dto.PaginatedRoleDTO{
		Items:      roleResponses,
		TotalItems: entity.TotalItems,
		TotalPages: totalPages,
		PageNumber: entity.PageNumber,
		PageSize:   entity.PageSize,
	}
}
