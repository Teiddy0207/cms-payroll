package service

import (
	"context"

	"cal-salary/core/cache"
	"cal-salary/core/errors"
	"cal-salary/core/params"
	"cal-salary/modules/auth/dto"
	"cal-salary/modules/auth/entity"
	"cal-salary/modules/auth/repository"

	"github.com/google/uuid"
)

type AuthService struct {
	repo  repository.AuthRepositoryInterface
	cache cache.Cache
	// storage storageRepo.StorageRepositoryInterface
}

func NewAuthService(repo repository.AuthRepositoryInterface, cache cache.Cache) AuthServiceInterface {
	return &AuthService{repo: repo, cache: cache}
}

type AuthServiceInterface interface {
	Register(ctx context.Context, requestData *dto.RegisterRequest) (*dto.RegisterResponse, *errors.AppError)
	Login(ctx context.Context, requestData *dto.LoginRequest) (*dto.LoginResponse, *errors.AppError)
	Logout(ctx context.Context, token string) *errors.AppError
	ChangePassword(ctx context.Context, token string, requestData *dto.ChangePasswordRequest) *errors.AppError
	ForgotPassword(ctx context.Context, identifier string) (*dto.ForgotPasswordResponse, *errors.AppError)
	VerifyOTP(ctx context.Context, requestData *dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, *errors.AppError)
	ResetPassword(ctx context.Context, requestData *dto.ResetPasswordRequest) *errors.AppError
	SendOTPChangePassword(ctx context.Context, token string) *errors.AppError

	PrivateCreateRole(ctx context.Context, role *dto.RoleRequest) (*dto.RoleResponse, *errors.AppError)
	PrivateGetRoles(ctx context.Context, params params.QueryParams) (*dto.PaginatedRoleDTO, *errors.AppError)
	PrivateGetRoleByID(ctx context.Context, id uuid.UUID) (*dto.RoleResponse, *errors.AppError)
	PrivateUpdateRole(ctx context.Context, id uuid.UUID, role *dto.RoleRequest) *errors.AppError
	PrivateDeleteRole(ctx context.Context, id uuid.UUID) *errors.AppError

	PrivateCreatePermission(ctx context.Context, permission *dto.PermissionRequest) *errors.AppError
	PrivateGetPermissions(ctx context.Context, params params.QueryParams) (*dto.PaginatedPermissionDTO, *errors.AppError)
	PrivateGetPermissionByID(ctx context.Context, id uuid.UUID) (*dto.PermissionResponse, *errors.AppError)
	PrivateUpdatePermission(ctx context.Context, id uuid.UUID, permission *dto.PermissionRequest) *errors.AppError
	PrivateDeletePermission(ctx context.Context, id uuid.UUID) *errors.AppError
	GetUserIDsByPermissionAndDepartment(ctx context.Context, permission string, departmentID uuid.UUID) ([]uuid.UUID, *errors.AppError)

	RefreshToken(ctx context.Context, token string) (*dto.RefreshTokenResponse, *errors.AppError)
	GetUserByIdentifier(ctx context.Context, identifier string) (*dto.UserResponse, *errors.AppError)

	PrivateAssignRoleToUser(ctx context.Context, req *dto.UserRoleRequest) *errors.AppError
	PrivateAssignPermissionToRole(ctx context.Context, req *dto.RolePermissionRequest) *errors.AppError
	PrivateAssignPermissionToUser(ctx context.Context, req *dto.UserPermissionRequest) *errors.AppError

	PrivateCreateUser(ctx context.Context, user *dto.CreateUserRequest) (*dto.UserDetailDTO, *errors.AppError)
	PrivateGetUserPermissions(ctx context.Context, userID uuid.UUID) ([]entity.Permission, *errors.AppError)
	PrivateGetUsers(ctx context.Context, params params.QueryParams) (*dto.PaginatedUserDTO, *errors.AppError)
	PrivateGetUser(ctx context.Context, userID uuid.UUID) (*dto.UserDetailDTO, *errors.AppError)
	PrivateGetPermissionsByUserID(ctx context.Context, userID uuid.UUID) (*[]dto.PermissionResponse, *errors.AppError)
	PrivateGetPermissionsByUserIDFromCache(ctx context.Context, userID uuid.UUID) (*[]dto.PermissionResponse, *errors.AppError)

	PrivateDeleteUser(ctx context.Context, id uuid.UUID) *errors.AppError
	PrivateUpdateUser(ctx context.Context, id uuid.UUID, user *dto.UserUpdateRequest) *errors.AppError
}
