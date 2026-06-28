package service

import (
	"cal-salary/core/constants"
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/entity"
	"cal-salary/modules/auth/mapper"
	"context"

	"github.com/google/uuid"
)

func (service *AuthService) PrivateAssignPermissionToUser(ctx context.Context, req *dto.UserPermissionRequest) *errors.AppError {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	err := service.repo.PrivateAssignPermissionToUser(ctx, mapper.ToUserPermissionEntity(req))
	if err != nil {
		logger.Error("AuthService:PrivateAssignPermissionToUser error: %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "Assign permission to user failed", err)
	}

	// Xóa cache permissions của user sau khi assign permission để đảm bảo permissions mới được load
	service.PrivateInvalidateUserPermissionsCache(ctx, req.UserID)

	return nil
}

func (service *AuthService) PrivateGetUserPermissions(ctx context.Context, userID uuid.UUID) ([]entity.Permission, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	result, err := service.repo.GetUserPermissions(ctx, userID)
	if err != nil {
		logger.Error("AuthService:GetUserPermissions:Error:", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "Get user permissions failed", err)
	}

	return result, nil
}
