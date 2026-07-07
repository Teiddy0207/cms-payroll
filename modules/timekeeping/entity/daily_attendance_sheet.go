package entity

import (
	"cal-salary/core/entity"
	"time"

	"github.com/google/uuid"
)

type DailyAttendanceSheet struct {
	EmployeeID    uuid.UUID  `db:"employee_id" json:"employee_id"`
	Date          time.Time  `db:"date" json:"date"`
	CheckIn       *time.Time `db:"check_in" json:"check_in"`
	CheckOut      *time.Time `db:"check_out" json:"check_out"`
	ActualWorkDay float64    `db:"actual_work_day" json:"actual_work_day"`
	OTHours       float64    `db:"ot_hours" json:"ot_hours"`
	Status        string     `db:"status" json:"status"`
	entity.BaseEntity
}
