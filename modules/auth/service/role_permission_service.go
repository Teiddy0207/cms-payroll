package service

import (
	"cal-salary/core/constants"
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/utils"
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/mapper"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (service *AuthService) PrivateGetPermissionsByUserIDFromCache(ctx context.Context, userID uuid.UUID) (*[]dto.PermissionResponse, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	key := utils.GenerateUserPermissionsKey(userID.String())

	val, err := service.cache.Get(ctx, key).Result()
	if err != nil {
		// Nếu key không tồn tại (redis.Nil), đây không phải là lỗi, chỉ là cache miss
		if err == redis.Nil {
			return nil, nil
		}
		logger.Error("AuthService:PrivateGetPermissionsByUserIDFromCache error: %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to get user permissions from cache", err)
	}

	permissionResponses := &[]dto.PermissionResponse{}
	errUnmarshal := json.Unmarshal([]byte(val), permissionResponses)
	if errUnmarshal != nil {
		logger.Error("AuthService:PrivateGetPermissionsByUserIDFromCache error: %v", errUnmarshal)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to unmarshal user permissions from cache", errUnmarshal)
	}

	return permissionResponses, nil
}

func (service *AuthService) PrivateGetPermissionsByUserID(ctx context.Context, userID uuid.UUID) (*[]dto.PermissionResponse, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultRequestTimeout)
	defer cancel()

	permissions, err := service.repo.PrivateGetPermissionsByUserID(ctx, userID)
	if err != nil {
		logger.Error("AuthService:PrivateGetPermissionsByUserID error: %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to get user permissions", err)
	}

	rolePermissions := mapper.ToPermissionDTOs(permissions)

	// Cache the rolePermissions
	key := utils.GenerateUserPermissionsKey(userID.String())
	rolePermissionsJSON, err := json.Marshal(rolePermissions)
	if err != nil {
		logger.Error("AuthService:PrivateGetPermissionsByUserID error: %v", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to marshal user permissions to cache", err)
	}
	errSet := service.cache.Set(ctx, key, rolePermissionsJSON, constants.DefaultZeroExpiration)
	if errSet != nil {
		logger.Error("AuthService:PrivateGetPermissionsByRole error: %v", errSet)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to set user permissions to cache", errSet)
	}

	return rolePermissions, nil
}

// PrivateInvalidateUserPermissionsCache xóa cache permissions của user
func (service *AuthService) PrivateInvalidateUserPermissionsCache(ctx context.Context, userID uuid.UUID) *errors.AppError {
	key := utils.GenerateUserPermissionsKey(userID.String())
	err := service.cache.Del(ctx, key)
	if err != nil {
		logger.Error("AuthService:PrivateInvalidateUserPermissionsCache error: %v", err)
		// Không trả về error vì việc xóa cache không nên làm fail operation chính
		return nil
	}
	logger.Info(fmt.Sprintf("AuthService:PrivateInvalidateUserPermissionsCache: Cache invalidated for user %s", userID.String()))
	return nil
}

func (service *AuthService) PrivateAssignPermissionToRole(ctx context.Context, req *dto.RolePermissionRequest) *errors.AppError {
	ctx, cancel := context.WithTimeout(ctx, constants.LongTimeout)
	defer cancel()

	req.GrantedBy = req.UserID

	err := service.repo.PrivateAssignPermissionToRole(ctx, req.RoleID, req.PermissionID, req.GrantedBy)
	if err != nil {
		logger.Error("AuthService:PrivateAssignPermissionToRole error: %v", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to assign permission to role", err)
	}

	// Xóa cache permissions của tất cả users có role này vì permissions của role đã thay đổi
	userIDs, err := service.repo.GetUserIDsByRoleID(ctx, req.RoleID)
	if err != nil {
		logger.Error("AuthService:PrivateAssignPermissionToRole:GetUserIDsByRoleID error: %v", err)
		// Không fail operation vì việc xóa cache là secondary
	} else {
		for _, userID := range userIDs {
			service.PrivateInvalidateUserPermissionsCache(ctx, userID)
		}
		logger.Info(fmt.Sprintf("AuthService:PrivateAssignPermissionToRole: Invalidated cache for %d users with role %s", len(userIDs), req.RoleID.String()))
	}

	return nil
}
