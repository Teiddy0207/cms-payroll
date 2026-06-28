package entity

import (
	"cal-salary/core/entity"

	"github.com/google/uuid"
)

type CompetencyEvaluation struct {
	EmployeeID       uuid.UUID `db:"employee_id" json:"employee_id"`
	EvaluatorID      uuid.UUID `db:"evaluator_id" json:"evaluator_id"`
	EvaluationPeriod string    `db:"evaluation_period" json:"evaluation_period"` // e.g. 2026-H1
	CompetencyID     uuid.UUID `db:"competency_id" json:"competency_id"`
	Score            int       `db:"score" json:"score"`
	Weight           int       `db:"weight" json:"weight"`
	entity.BaseEntity
}
