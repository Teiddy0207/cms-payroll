package entity

import (
	"cal-salary/core/entity"

	"github.com/google/uuid"
)

type Department struct {
	Code        string     `db:"code" json:"code"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	ParentID    *uuid.UUID `db:"parent_id" json:"parent_id"`
	ManagerID   *uuid.UUID `db:"manager_id" json:"manager_id"`
	entity.BaseEntity
}

type PaginatedDepartmentEntity = entity.Pagination[Department]
