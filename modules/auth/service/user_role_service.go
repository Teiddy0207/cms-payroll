package service

import (
	"context"

	"cal-salary/core/constants"
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/mapper"
)

func (service *AuthService) PrivateAssignRoleToUser(ctx context.Context, req *dto.UserRoleRequest) *errors.AppError {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	err := service.repo.PrivateAssignRoleToUser(ctx, mapper.ToUserRoleEntity(req))
	if err != nil {
		logger.Error("AuthService:PrivateAssignRoleToUser error: %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "Assign role to user failed", err)
	}

	// Xóa cache permissions của user sau khi assign role để đảm bảo permissions mới được load
	service.PrivateInvalidateUserPermissionsCache(ctx, req.UserID)

	return nil
}
