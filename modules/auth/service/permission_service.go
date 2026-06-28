package service

import (
	"cal-salary/core/constants"
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/mapper"
	"context"

	"github.com/google/uuid"
)

func (service *AuthService) PrivateCreatePermission(ctx context.Context, permission *dto.PermissionRequest) *errors.AppError {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	err := service.repo.PrivateCreatePermission(ctx, mapper.ToPermissionEntity(permission))
	if err != nil {
		logger.Error("AuthService:PrivateCreatePermission error: %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "Create permission failed", err)
	}

	return nil
}

func (service *AuthService) PrivateGetPermissions(ctx context.Context, params params.QueryParams) (*dto.PaginatedPermissionDTO, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	permissions, err := service.repo.PrivateGetPermissions(ctx, params)
	if err != nil {
		logger.Error("AuthService:PrivateGetPermissions error: %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "Get permissions failed", err)
	}

	return mapper.ToPaginatedPermissionDTO(permissions), nil
}

func (service *AuthService) PrivateGetPermissionByID(ctx context.Context, id uuid.UUID) (*dto.PermissionResponse, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	permission, err := service.repo.PrivateGetPermissionByID(ctx, id)
	if err != nil {
		logger.Error("AuthService:PrivateGetPermissionByID error: %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "Get permission failed", err)
	}

	return mapper.ToPermissionDTO(permission), nil
}

func (service *AuthService) PrivateUpdatePermission(ctx context.Context, id uuid.UUID, permission *dto.PermissionRequest) *errors.AppError {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	err := service.repo.PrivateUpdatePermission(ctx, id, mapper.ToPermissionEntity(permission))
	if err != nil {
		logger.Error("AuthService:PrivateUpdatePermission error: %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "Update permission failed", err)
	}

	return nil
}

func (service *AuthService) PrivateDeletePermission(ctx context.Context, id uuid.UUID) *errors.AppError {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	err := service.repo.PrivateDeletePermission(ctx, id)
	if err != nil {
		logger.Error("AuthService:PrivateDeletePermission error: %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "Delete permission failed", err)
	}

	return nil
}
