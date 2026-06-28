package entity

import (
	"cal-salary/core/entity"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	Name        string  `db:"name"`
	Slug        string  `db:"slug"`
	Description *string `db:"description"`
	IsSystem    bool    `db:"is_system"`
	IsActive    bool    `db:"is_active"`
	entity.BaseEntity
}

type UserRole struct {
	ID          string     `db:"id"`
	Code        string     `db:"code"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`
	UserID      uuid.UUID  `db:"user_id"`
	RoleID      uuid.UUID  `db:"role_id"`
	AssignedBy  uuid.UUID  `db:"assigned_by"`
	AssignedAt  time.Time  `db:"assigned_at"`
	ExpiresAt   *time.Time `db:"expires_at"`
	IsActive    bool       `db:"is_active"`
	entity.BaseEntity
}

type PaginatedRoleEntity = entity.Pagination[Role]
