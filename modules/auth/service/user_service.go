package service

import (
	"cal-salary/core/constants"
	"cal-salary/core/errors"
	"cal-salary/core/logger"
	"cal-salary/core/params"
	"cal-salary/core/utils"
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/mapper"
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (service *AuthService) PrivateCreateUser(ctx context.Context, user *dto.CreateUserRequest) (*dto.UserDetailDTO, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	entity := mapper.ToUserEntity(user)
	entity.Password, _ = utils.HashPassword(entity.Password)

	createdUser, err := service.repo.PrivateCreateUser(ctx, entity)
	if err != nil {
		logger.Error("AuthService:PrivateCreateUser:Error:", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to create user", err)
	}

	if user.RoleID != nil && *user.RoleID != uuid.Nil {
		userRoleReq := &dto.UserRoleRequest{
			UserID: createdUser.ID,
			// RoleID:     *user.RoleID,
			// AssignedAt: createdUser.CreatedAt,
			// IsActive:   true,
		}

		err = service.PrivateAssignRoleToUser(ctx, userRoleReq)
		if err != nil {
			logger.Error("AuthService:PrivateCreateUser:AssignRole:Error:", err)
		}
	}

	// Lấy lại user với đầy đủ thông tin từ repository (không qua service để tránh tạo context mới)
	userDetailEntity, err := service.repo.PrivateGetUser(ctx, createdUser.ID)
	if err != nil {
		logger.Error("AuthService:PrivateCreateUser:GetUser:Error:", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to get created user", err)
	}
	if userDetailEntity == nil {
		logger.Error("AuthService:PrivateCreateUser:UserNotFound:", createdUser.ID)
		return nil, errors.NewAppError(errors.ErrInternalServer, "created user not found", nil)
	}

	return mapper.ToUserDetailDTO(userDetailEntity), nil
}

func (service *AuthService) PrivateGetUser(ctx context.Context, userID uuid.UUID) (*dto.UserDetailDTO, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	user, err := service.repo.PrivateGetUser(ctx, userID)
	if err != nil {
		logger.Error("AuthService:PrivateGetUser:Error:", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to get user", err)
	}

	if user == nil {
		fmt.Println("AuthService:PrivateGetUser:UserNotFound:", userID)
		return nil, errors.NewAppError(errors.ErrNotFound, "user not found", nil)
	}

	return mapper.ToUserDetailDTO(user), nil
}

func (service *AuthService) PrivateGetUsers(ctx context.Context, params params.QueryParams) (*dto.PaginatedUserDTO, *errors.AppError) {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	users, err := service.repo.PrivateGetUsers(ctx, params)
	if err != nil {
		logger.Error("AuthService:PrivateGetUsers:Error:", err)
		return nil, errors.NewAppError(errors.ErrInternalServer, "failed to get users", err)
	}

	return mapper.ToUserWithProfilePaginationDTO(users), nil
}

func (service *AuthService) PrivateDeleteUser(ctx context.Context, id uuid.UUID) *errors.AppError {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()
	existingUser, err := service.repo.PrivateGetUser(ctx, id)
	if err != nil {
		logger.Error("AuthService:PrivateDeleteUser", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to get user by id", err)
	}
	if existingUser == nil {
		fmt.Println("AuthService:PrivateDeleteUser:userNotFound")
		return errors.NewAppError(errors.ErrNotFound, "user not found", nil)
	}
	errDelete := service.repo.PrivateDeleteUser(ctx, id)
	if errDelete != nil {
		logger.Error("AuthService:PrivateDeleteUser", errDelete)
		return errors.NewAppError(errors.ErrInternalServer, "failed to delete user", errDelete)
	}
	return nil

}

func (service *AuthService) PrivateUpdateUser(ctx context.Context, id uuid.UUID, user *dto.UserUpdateRequest) *errors.AppError {
	ctx, cancel := context.WithTimeout(ctx, constants.DefaultTimeout)
	defer cancel()

	// Lấy existing user entity từ identifier (có thể dùng UUID string)
	existingUser, err := service.repo.GetUserByIdentifier(ctx, id.String())
	if err != nil {
		logger.Error("AuthService:PrivateUpdateUser:GetUserByIdentifier:Error:", err)
		return errors.NewAppError(errors.ErrInternalServer, "failed to get user", err)
	}

	if existingUser == nil {
		return errors.NewAppError(errors.ErrNotFound, "user not found", nil)
	}

	// Map DTO sang entity với existing data
	entity := mapper.ToUserUpdateEntity(user, existingUser)
	entity.ID = id

	// Hash password nếu có
	if user.Password != nil && *user.Password != "" {
		hashedPassword, err := utils.HashPassword(*user.Password)
		if err != nil {
			logger.Error("AuthService:PrivateUpdateUser:HashPassword:Error:", err)
			return errors.NewAppError(errors.ErrInternalServer, "failed to hash password", err)
		}
		entity.Password = hashedPassword
	}

	err = service.repo.PrivateUpdateUser(ctx, entity, id)
	if err != nil {
		logger.Error("AuthService:PrivateUpdateUser:Error:", err)
		if err.Error() == "user not found" {
			return errors.NewAppError(errors.ErrNotFound, "user not found", err)
		}
		return errors.NewAppError(errors.ErrInternalServer, "failed to update user", err)
	}

	return nil
}
