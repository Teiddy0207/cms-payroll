package entity

import (
	"cal-salary/core/entity"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Email           *string    `db:"email"`
	Username        *string    `db:"username"`
	Password        string     `db:"password"`
	EmailVerifiedAt *time.Time `db:"email_verified_at"`
	PositionID      *string    `db:"position_id"`
	IsActive        bool       `db:"is_active"`
	entity.BaseEntity
}

type UserDetail struct {
	ID          string     `db:"id"`
	Email       *string    `db:"email"`
	Phone       string     `db:"phone"`
	Username    *string    `db:"username"`
	IsActive    bool       `db:"is_active"`
	CreatedAt   string     `db:"created_at"`
	DisplayName *string    `db:"display_name"`
	FullName    *string    `db:"full_name"`
	Avatar      *string    `db:"avatar"`
	DateOfBirth *string    `db:"date_of_birth"`
	Gender      *string    `db:"gender"`
	RoleID      *uuid.UUID `db:"role_id"`
	RoleCode    *string    `db:"role_code"`
	RoleName    *string    `db:"role_name"`
	RoleDesc    *string    `db:"role_description"`
	Permissions *string    `db:"permissions"` // JSON string chứa danh sách permissions
}

type UserWithProfile struct {
	User
	ProfileID       *uuid.UUID `db:"profile_id"`
	ProfileCode     *string    `db:"profile_code"`
	ProfileFullName *string    `db:"profile_full_name"`
	ProfilePhone    *string    `db:"profile_phone"`
	ProfileGender   *string    `db:"profile_gender"`
	ProfileAvatar   *string    `db:"profile_avatar"`
	OfficialID      *string    `db:"official_id"`
	OfficeCode      *string    `db:"office_code"`
	OfficeName      *string    `db:"office_name"`
	RoleID          *uuid.UUID `db:"role_id"`
	RoleCode        *string    `db:"role_code"`
	RoleName        *string    `db:"role_name"`
	RoleDesc        *string    `db:"role_description"`
	PositionID      *uuid.UUID `db:"position_id"`
	PositionCode    *string    `db:"position_code"`
	PositionName    *string    `db:"position_name"`
	Roles           *string    `db:"roles"` // JSON string chứa danh sách roles
}


type PaginatedUserEntity = entity.Pagination[User]
type PaginatedUserWithProfileEntity = entity.Pagination[UserWithProfile]
