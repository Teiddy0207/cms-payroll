package entity

import (
	"cal-salary/core/entity"
	"time"

	"github.com/google/uuid"
)

type ExplanationRequest struct {
	EmployeeID uuid.UUID  `db:"employee_id" json:"employee_id"`
	Date       time.Time  `db:"date" json:"date"`
	Reason     string     `db:"reason" json:"reason"`
	ApprovedBy *uuid.UUID `db:"approved_by" json:"approved_by"`
	Status     string     `db:"status" json:"status"`
	entity.BaseEntity
}
