package service

import (
	"cal-salary/core/constants"
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/entity"
	"cal-salary/modules/auth/mapper"
	"context"

	"github.com/google/uuid"
)

func (service *AuthService) PrivateCreateRole(ctx context.Context, role *dto.RoleRequest) (*dto.RoleResponse, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	roleEntity := mapper.ToRoleEntity(role)
	createdRole, err := service.repo.PrivateCreateRole(ctx, roleEntity)
	if err != nil {
		logger.Error("AuthService:PrivateCreateRole error: %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "Create role failed", err)
	}

	return mapper.ToRoleDTO(createdRole, nil), nil
}

func (service *AuthService) PrivateGetRoles(ctx context.Context, params params.QueryParams) (*dto.PaginatedRoleDTO, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	roles, err := service.repo.PrivateGetRoles(ctx, params)
	if err != nil {
		logger.Error("AuthService:PrivateGetRoles error: %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "Get roles failed", err)
	}

	rolePermissionsMap := make(map[uuid.UUID]*[]entity.Permission)
	if roles != nil && len(roles.Items) > 0 {
		for _, role := range roles.Items {
			permissions, err := service.repo.GetPermissionsByRoleID(ctx, role.ID)
			if err != nil {
				logger.Error("AuthService:PrivateGetRoles:GetPermissionsByRoleID error for role %s: %v", role.ID.String(), err)
				rolePermissionsMap[role.ID] = nil
			} else {
				rolePermissionsMap[role.ID] = permissions
			}
		}
	}

	return mapper.ToPaginatedRoleDTO(roles, rolePermissionsMap), nil
}

func (service *AuthService) PrivateGetRoleByID(ctx context.Context, id uuid.UUID) (*dto.RoleResponse, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	role, err := service.repo.PrivateGetRoleByID(ctx, id)
	if err != nil {
		logger.Error("AuthService:PrivateGetRoleByID error: %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "Get role by ID failed", err)
	}

	if role == nil {
		return nil, errors.NewAppError(errors.ErrNotFound, "Role not found", nil)
	}

	// Lấy permissions của role
	permissions, err := service.repo.GetPermissionsByRoleID(ctx, id)
	if err != nil {
		logger.Error("AuthService:PrivateGetRoleByID:GetPermissionsByRoleID error: %v", err)
		// Nếu lỗi, set permissions là nil
		permissions = nil
	}

	return mapper.ToRoleDTO(role, permissions), nil
}

func (service *AuthService) PrivateUpdateRole(ctx context.Context, id uuid.UUID, role *dto.RoleRequest) *errors.AppError {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	err := service.repo.PrivateUpdateRole(ctx, id, mapper.ToRoleEntity(role))
	if err != nil {
		logger.Error("AuthService:PrivateUpdateRole error: %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "Update role failed", err)
	}

	return nil
}

func (service *AuthService) PrivateDeleteRole(ctx context.Context, id uuid.UUID) *errors.AppError {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	err := service.repo.PrivateDeleteRole(ctx, id)
	if err != nil {
		logger.Error("AuthService:PrivateDeleteRole error: %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "Delete role failed", err)
	}

	return nil
}
