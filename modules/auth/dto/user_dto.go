package dto

import (
	"cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username          string     `json:"username"`
	Password          string     `json:"password"`
	ConfirmedPassword string     `json:"confirmed_password"`
	RoleID            *uuid.UUID `json:"role_id,omitempty"`
}

type RegisterRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Identifier string `json:"identifiers"` // username, email
	Password   string `json:"password"`
}

type LoginUserProfile struct {
	ID       uuid.UUID     `json:"id"`
	UserID   *uuid.UUID    `json:"user_id"`
	User     *UserInfo     `json:"user,omitempty"`
	FullName *string       `json:"full_name"`
	Position *PositionInfo `json:"position,omitempty"`
}

type PositionInfo struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type UserInfo struct {
	ID       *uuid.UUID `json:"id"`
	Username *string    `json:"username"`
	Email    *string    `json:"email"`
}

type UserProfileResponse struct {
	ID uuid.UUID `json:"id"`
}

type LoginResponse struct {
	AccessToken   string            `json:"access_token"`
	RefreshToken  string            `json:"refresh_token"`
	UserProfileID *uuid.UUID        `json:"user_profile_id"`
	UserProfile   *LoginUserProfile `json:"user_profile,omitempty"`
}

type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

type ChangePasswordRequest struct {
	dto.BaseRequest
	Password        string `json:"password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
	OTP             string `json:"otp"`
}

type ForgotPasswordRequest struct {
	Identifier string `json:"identifier"`
}

type ForgotPasswordResponse struct {
	UserId uuid.UUID `json:"user_id"`
}

type VerifyOTPRequest struct {
	UserID uuid.UUID `json:"user_id"`
	OTP    string    `json:"otp"`
}

type VerifyOTPResponse struct {
	Token string `json:"token"`
}

type ResetPassword struct {
	Token             string `json:"token"`
	NewPassword       string `json:"new_password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type UserRequest struct {
	Email           string     `json:"email"`
	Phone           string     `json:"phone"`
	Username        *string    `json:"username"`
	Password        *string    `json:"-"`
	EmailVerifiedAt *time.Time `json:"-"`
	PhoneVerifiedAt *time.Time `json:"-"`
	LockedUntil     *time.Time `json:"-"`
	IsActive        bool       `json:"is_active"`
}

type UserUpdateRequest struct {
	Email      *string `json:"email,omitempty"`
	Username   *string `json:"username,omitempty"`
	Password   *string `json:"password,omitempty"`
	PositionID *string `json:"position_id,omitempty"`
	IsActive   *bool   `json:"is_active,omitempty"`
}

type UserProfileInfo struct {
	ID       *uuid.UUID `json:"id,omitempty"`
	Code     *string    `json:"code,omitempty"`
	FullName *string    `json:"full_name,omitempty"`
	Phone    *string    `json:"phone,omitempty"`
	Gender   *string    `json:"gender,omitempty"`
	Avatar   *string    `json:"avatar,omitempty"`
}

type OfficeInfo struct {
	Code *string `json:"code,omitempty"`
	Name *string `json:"name,omitempty"`
}

type RoleInfo struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
}

type UserResponse struct {
	ID              uuid.UUID        `json:"id"`
	Email           *string          `json:"email"`
	Username        *string          `json:"username"`
	Password        string           `json:"-"`
	EmailVerifiedAt *time.Time       `json:"email_verified_at"`
	PhoneVerifiedAt *time.Time       `json:"phone_verified_at"`
	LockedUntil     *time.Time       `json:"locked_until"`
	IsActive        bool             `json:"is_active"`
	UserProfile     *UserProfileInfo `json:"user_profile,omitempty"`
	OfficialID      *string          `json:"official_id,omitempty"`
	Official        *OfficeInfo      `json:"official,omitempty"`
	PositionID      *uuid.UUID       `json:"position_id,omitempty"`
	Position        *PositionInfo    `json:"position,omitempty"`
	Roles           *RoleInfo        `json:"roles,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type UserDetailDTO struct {
	ID          string               `json:"id"`
	Email       *string              `json:"email"`
	Phone       string               `json:"phone"`
	Username    *string              `json:"username"`
	IsActive    bool                 `json:"is_active"`
	CreatedAt   string               `json:"created_at"`
	DisplayName *string              `json:"display_name"`
	FullName    *string              `json:"full_name"`
	Avatar      *string              `json:"avatar"`
	DateOfBirth *string              `json:"date_of_birth"`
	Gender      *string              `json:"gender"`
	Roles       *RoleInfo            `json:"roles,omitempty"`
	Permissions []PermissionResponse `json:"permissions,omitempty"`
}

type PaginatedUserDTO = dto.Pagination[UserResponse]
