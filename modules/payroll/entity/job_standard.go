package entity

import (
	"cal-salary/core/entity"
)

type JobStandard struct {
	StandardCode   string  `db:"standard_code" json:"standard_code"`
	Name           string  `db:"name" json:"name"`
	AllowanceValue float64 `db:"allowance_value" json:"allowance_value"`
	Description    *string `db:"description" json:"description"`
	entity.BaseEntity
}
