package dto

import (
	"cal-salary/core/dto"
	"time"

	"github.com/google/uuid"
)

type CreateUserProfileRequest struct {
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	Code         string     `json:"code" validate:"required,min=2,max=50"`
	FullName     string     `json:"full_name" validate:"required,min=2,max=255"`
	Phone        *string    `json:"phone,omitempty"`
	Avatar       *string    `json:"avatar,omitempty"`
	DateOfBirth  *string    `json:"date_of_birth,omitempty"`
	Gender       *string    `json:"gender,omitempty"`
	PositionID   *uuid.UUID `json:"position_id,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
}

type UpdateUserProfileRequest struct {
	FullName     string     `json:"full_name" validate:"required,min=2,max=255"`
	Phone        *string    `json:"phone,omitempty"`
	Avatar       *string    `json:"avatar,omitempty"`
	DateOfBirth  *string    `json:"date_of_birth,omitempty"`
	Gender       *string    `json:"gender,omitempty"`
	PositionID   *uuid.UUID `json:"position_id,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
}

type UserProfileResponse struct {
	ID           uuid.UUID  `json:"id"`
	UserID       *uuid.UUID `json:"user_id,omitempty"`
	Code         string     `json:"code"`
	FullName     string     `json:"full_name"`
	Phone        *string    `json:"phone,omitempty"`
	Avatar       *string    `json:"avatar,omitempty"`
	DateOfBirth  *string    `json:"date_of_birth,omitempty"`
	Gender       *string    `json:"gender,omitempty"`
	PositionID   *uuid.UUID `json:"position_id,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type PaginatedUserProfileDTO = dto.Pagination[UserProfileResponse]

