package entity

import (
	"cal-salary/core/entity"

	"github.com/google/uuid"
)

type JobPosition struct {
	Code         string     `db:"code" json:"code"`
	Name         string     `db:"name" json:"name"`
	Description  *string    `db:"description" json:"description"`
	DepartmentID *uuid.UUID `db:"department_id" json:"department_id"`
	entity.BaseEntity
}
type PaginatedJobPositionEntity = entity.Pagination[JobPosition]
