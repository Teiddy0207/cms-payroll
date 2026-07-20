package entity

import (
	"time"

	"github.com/google/uuid"
)

type PayrollFormula struct {
	ID           uuid.UUID  `db:"id" json:"id"`
	VariableName string     `db:"variable_name" json:"variable_name"`
	Expression   string     `db:"expression" json:"expression"`
	StartDate    time.Time  `db:"start_date" json:"start_date"`
	EndDate      *time.Time `db:"end_date" json:"end_date"`
	Description  string     `db:"description" json:"description"`
	CreatedAt    time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updated_at"`
}
