package entity

import (
	"cal-salary/core/entity"
	"time"

	"github.com/google/uuid"
)

type OTRequest struct {
	EmployeeID     uuid.UUID  `db:"employee_id" json:"employee_id"`
	Date           time.Time  `db:"date" json:"date"`
	HoursRequested float64    `db:"hours_requested" json:"hours_requested"`
	IsNightOT      bool       `db:"is_night_ot" json:"is_night_ot"`
	IsHolidayOT    bool       `db:"is_holiday_ot" json:"is_holiday_ot"`
	ApprovedBy     *uuid.UUID `db:"approved_by" json:"approved_by"`
	Status         string     `db:"status" json:"status"`
	entity.BaseEntity
}
