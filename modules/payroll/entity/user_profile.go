package entity

import (
	"cal-salary/core/entity"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserID       *uuid.UUID `db:"user_id" json:"user_id"`
	Code         string     `db:"code" json:"code"`
	FullName     string     `db:"full_name" json:"full_name"`
	Phone        *string    `db:"phone" json:"phone"`
	Avatar       *string    `db:"avatar" json:"avatar"`
	DateOfBirth  *string    `db:"date_of_birth" json:"date_of_birth"`
	Gender       *string    `db:"gender" json:"gender"`
	PositionID   *uuid.UUID `db:"position_id" json:"position_id"`
	DepartmentID *uuid.UUID `db:"department_id" json:"department_id"`

	JobPosition *JobPosition `db:"-" json:"job_position,omitempty"`
	entity.BaseEntity
}

type PaginatedUserProfileEntity = entity.Pagination[UserProfile]
