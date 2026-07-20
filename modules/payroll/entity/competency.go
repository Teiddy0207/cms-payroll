package entity

import (
	"cal-salary/core/entity"

	"github.com/google/uuid"
)

type Competency struct {
	Code        string     `db:"code" json:"code"`
	Name        string     `db:"name" json:"name"`
	Description *string    `db:"description" json:"description"`
	TypeID      *uuid.UUID `db:"type_id" json:"type_id"`
	PointValue  int        `db:"point_value" json:"point_value"`
	entity.BaseEntity
}

type PaginatedCompetencyEntity = entity.Pagination[Competency]
