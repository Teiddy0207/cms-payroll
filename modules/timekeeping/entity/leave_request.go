package entity

import (
	"cal-salary/core/entity"
	"time"

	"github.com/google/uuid"
)

// LeaveRequest đại diện cho đơn xin nghỉ phép của nhân viên.
// DaysRequested hỗ trợ 0.5 (nửa ngày) hoặc số nguyên (cả ngày).
// Khi start_date == end_date và days_requested == 0.5 → đây là đơn nửa ngày.
type LeaveRequest struct {
	EmployeeID    uuid.UUID  `db:"employee_id"    json:"employee_id"`
	StartDate     time.Time  `db:"start_date"     json:"start_date"`
	EndDate       time.Time  `db:"end_date"       json:"end_date"`
	DaysRequested float64    `db:"days_requested" json:"days_requested"`
	Reason        string     `db:"reason"         json:"reason"`
	ApprovedBy    *uuid.UUID `db:"approved_by"    json:"approved_by"`
	Status        string     `db:"status"         json:"status"` // PENDING, APPROVED, REJECTED
	entity.BaseEntity
}
